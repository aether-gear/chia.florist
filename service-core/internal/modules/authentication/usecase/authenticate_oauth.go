package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	customerDomain "service-core/internal/modules/customer/domain"
	customerRepo "service-core/internal/modules/customer/repository"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type AuthenticateOAuthUsecase struct {
	executor         transaction.Executor
	transactor       transaction.Transactor
	accountRepo      repository.AccountRepository
	oauthRepo        repository.OAuthConnectionRepository
	userRepo         userRepo.UserRepository
	customerRepo     customerRepo.CustomerRepository
	tokenHasher      repository.TokenHasher
	tokenSvc         repository.TokenService
	sessionRepo      repository.SessionRepository
	refreshTokenRepo repository.RefreshTokenRepository
	auditLogger      applogger.AuditLogger
}

func NewAuthenticateOAuthUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	oauthRepo repository.OAuthConnectionRepository,
	userRepo userRepo.UserRepository,
	customerRepo customerRepo.CustomerRepository,
	tokenHasher repository.TokenHasher,
	tokenSvc repository.TokenService,
	sessionRepo repository.SessionRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	auditLogger applogger.AuditLogger,
) *AuthenticateOAuthUsecase {
	return &AuthenticateOAuthUsecase{
		executor:         executor,
		transactor:       transactor,
		accountRepo:      accountRepo,
		oauthRepo:        oauthRepo,
		userRepo:         userRepo,
		customerRepo:     customerRepo,
		tokenHasher:      tokenHasher,
		tokenSvc:         tokenSvc,
		sessionRepo:      sessionRepo,
		refreshTokenRepo: refreshTokenRepo,
		auditLogger:      auditLogger,
	}
}

type AuthenticateOAuthParams struct {
	UserAgent *string
	IPAddress *string
	Provider  domain.OAuthProvider
	Subject   string
	Email     string
	Name      string
	AvatarURL *string
}

type AuthenticateOAuthResult struct {
	AccessToken  repository.GeneratedToken
	RefreshToken repository.GeneratedToken
}

func (u *AuthenticateOAuthUsecase) Execute(
	ctx context.Context,
	input AuthenticateOAuthParams,
) (*AuthenticateOAuthResult, error) {
	now := time.Now()

	conn, err := u.oauthRepo.
		GetByProviderAndSubject(
			ctx,
			u.executor,
			input.Provider,
			input.Subject,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve oauth connection: %w", err)
	}

	var userID uuid.UUID
	var account *domain.Account

	// OAuth identity is already linked;
	// update login metadata and continue authentication flow
	// with existing account.
	if conn != nil {
		userID = conn.UserID

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			if err := u.oauthRepo.
				UpdateLastLogin(ctx, exec, conn.ID, now); err != nil {
				return fmt.Errorf("failed to update last login: %w", err)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		account, err = u.accountRepo.
			GetByUserID(ctx, u.executor, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get account: %w", err)
		}

		if account == nil {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category: "user_action",
				Action:   "oauth_login",
				Resource: "session",
				Outcome:  applogger.OutcomeFailure,
				Metadata: map[string]any{
					"email":    input.Email,
					"provider": string(input.Provider),
					"reason":   "account not found",
				},
			})
			return nil, apperrors.NewNotFound("account not found")
		}

		if account.Status != domain.AccountActive &&
			account.Status != domain.AccountPending {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category: "user_action",
				Action:   "oauth_login",
				Resource: "session",
				Outcome:  applogger.OutcomeFailure,
				Metadata: map[string]any{
					"email":    input.Email,
					"provider": string(input.Provider),
					"reason":   "account suspended or locked",
				},
			})
			return nil, apperrors.NewForbidden("account is suspended or locked")
		}

		if account.Status == domain.AccountPending {
			err = u.accountRepo.
				ActivateByUserID(ctx, u.executor, userID)
			if err != nil {
				return nil, fmt.Errorf("failed to activate account: %w", err)
			}

			account.Status = domain.AccountActive
		}
	}

	// No linked OAuth identity found;
	// check whether the email belongs to an existing account.
	if conn == nil {
		account, err = u.accountRepo.
			GetByEmail(ctx, u.executor, input.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve account by email: %w", err)
		}
	}

	// No linked OAuth identity found;
	// check whether the email belongs to an existing account.
	if conn == nil && account != nil {
		// The schema enforces user_id UNIQUE,
		// so check if user already have any OAuth connection.
		existingConn, err := u.oauthRepo.
			GetByUserID(ctx, u.executor, account.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing connection by user id: %w", err)
		}

		if existingConn != nil {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category: "user_action",
				Action:   "oauth_login",
				Resource: "session",
				Outcome:  applogger.OutcomeFailure,
				Metadata: map[string]any{
					"email":    input.Email,
					"provider": string(input.Provider),
					"reason":   "email already linked to another OAuth account",
				},
			})
			return nil, apperrors.NewConflict("this email is already linked with another OAuth account")
		}

		userID = account.UserID
		newConn := domain.OAuthConnection{
			ID:          uuid.New(),
			UserID:      userID,
			Provider:    input.Provider,
			Subject:     input.Subject,
			Email:       &input.Email,
			LastLoginAt: &now,
			CreatedAt:   now,
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			if account.Status == domain.AccountPending {
				if err := u.accountRepo.
					ActivateByUserID(ctx, exec, userID); err != nil {
					return fmt.Errorf("failed to activate account: %w", err)
				}

				account.Status = domain.AccountActive
			}

			if err := u.oauthRepo.
				Create(ctx, exec, newConn); err != nil {
				return fmt.Errorf("failed to link oauth connection: %w", err)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// Neither OAuth identity nor local Account found;
	// Create a new User, Customer, Account, and OAuth Connection.
	if conn == nil && account == nil {
		userID = uuid.New()
		accountID := uuid.New()
		customerID := uuid.New()

		baseUsername, _, _ := strings.Cut(input.Email, "@")
		username := fmt.Sprintf("%s_%s", baseUsername, uuid.New().String()[:7])

		userProps := userRepo.CreateUserProps{
			ID:        userID,
			Name:      input.Name,
			Username:  username,
			AvatarURL: input.AvatarURL,
			CreatedAt: now,
		}

		newAcc := domain.Account{
			ID:        accountID,
			UserID:    userID,
			Email:     input.Email,
			Password:  "",
			Status:    domain.AccountActive,
			Type:      domain.AccountTypeCustomer,
			CreatedAt: now,
		}

		cust := customerDomain.Customer{
			ID:        customerID,
			UserID:    userID,
			CreatedAt: now,
		}

		newConn := domain.OAuthConnection{
			ID:          uuid.New(),
			UserID:      userID,
			Provider:    input.Provider,
			Subject:     input.Subject,
			Email:       &input.Email,
			LastLoginAt: &now,
			CreatedAt:   now,
		}

		profile := userRepo.SaveProfileProps{
			UserID:    userID,
			Name:      &input.Name,
			AvatarURL: input.AvatarURL,
			UpdatedAt: now,
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			if err := u.userRepo.
				CreateUser(ctx, exec, userProps); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}

			if err := u.userRepo.
				SaveProfile(ctx, exec, profile); err != nil {
				return fmt.Errorf("failed to save user profile: %w", err)
			}

			if err := u.customerRepo.
				Create(ctx, exec, cust); err != nil {
				return fmt.Errorf("failed to create customer profile: %w", err)
			}

			if err := u.accountRepo.
				Create(ctx, exec, newAcc); err != nil {
				return fmt.Errorf("failed to create account: %w", err)
			}

			if err := u.oauthRepo.
				Create(ctx, exec, newConn); err != nil {
				return fmt.Errorf("failed to create oauth connection: %w", err)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		account = &newAcc
	}

	session := domain.Session{
		ID:        uuid.New(),
		UserID:    userID,
		UserAgent: input.UserAgent,
		IPAddress: input.IPAddress,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	var customerID *uuid.UUID
	custProfile, err := u.customerRepo.
		GetByUserID(ctx, u.executor, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve customer profile: %w", err)
	}
	if custProfile != nil {
		customerID = &custProfile.ID
	}

	tknAccessInput := repository.GenerateTokenParams{
		UserID:     userID,
		SessionID:  session.ID,
		CustomerID: customerID,
		Type:       domain.TokenTypeAccess,
		Duration:   30 * time.Minute,
	}
	accessTkn, err := u.tokenSvc.Generate(tknAccessInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	tknRefreshInput := repository.GenerateTokenParams{
		UserID:     userID,
		SessionID:  session.ID,
		CustomerID: customerID,
		Type:       domain.TokenTypeRefresh,
		Duration:   7 * 24 * time.Hour,
	}
	refreshTkn, err := u.tokenSvc.Generate(tknRefreshInput)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	refreshTknHashed := u.tokenHasher.Hash(refreshTkn.Token)
	refreshTknDomain := domain.RefreshToken{
		ID:        uuid.New(),
		SessionID: session.ID,
		TokenHash: refreshTknHashed,
		ExpiresAt: now.Add(7 * 24 * time.Hour),
		CreatedAt: now,
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.sessionRepo.
			Save(ctx, exec, session); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		if err := u.refreshTokenRepo.
			Save(ctx, exec, refreshTknDomain); err != nil {
			return fmt.Errorf("failed to save refresh token: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category:   "user_action",
		Action:     "oauth_login",
		Resource:   "session",
		ResourceID: session.ID.String(),
		Outcome:    applogger.OutcomeSuccess,
		Metadata: map[string]any{
			"email":    input.Email,
			"provider": string(input.Provider),
		},
	})

	return &AuthenticateOAuthResult{
		AccessToken:  accessTkn,
		RefreshToken: refreshTkn,
	}, nil
}
