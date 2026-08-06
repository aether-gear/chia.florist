package usecase

import (
	"context"
	"fmt"
	"time"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	customerDomain "service-core/internal/modules/customer/domain"
	customerRepo "service-core/internal/modules/customer/repository"
	userRepo "service-core/internal/modules/user/repository"
	mailer "service-core/internal/shared/mailer"
	otp "service-core/internal/shared/otp"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type RegisterCustomerUsecase struct {
	executor      transaction.Executor
	transactor    transaction.Transactor
	accountRepo   repository.AccountRepository
	hasher        domain.PasswordHasher
	userRepo      userRepo.UserRepository
	customerRepo  customerRepo.CustomerRepository
	challengeRepo repository.VerificationChallengeRepository
	otpGen        otp.Generator
	mailer        mailer.Sender
	auditLogger   applogger.AuditLogger
}

func NewRegisterCustomerUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo repository.AccountRepository,
	hasher domain.PasswordHasher,
	userRepo userRepo.UserRepository,
	customerRepo customerRepo.CustomerRepository,
	challengeRepo repository.VerificationChallengeRepository,
	otpGen otp.Generator,
	mailer mailer.Sender,
	auditLogger applogger.AuditLogger,
) *RegisterCustomerUsecase {
	return &RegisterCustomerUsecase{
		executor:      executor,
		transactor:    transactor,
		accountRepo:   accountRepo,
		hasher:        hasher,
		userRepo:      userRepo,
		customerRepo:  customerRepo,
		challengeRepo: challengeRepo,
		otpGen:        otpGen,
		mailer:        mailer,
		auditLogger:   auditLogger,
	}
}

type RegisterCustomerParams struct {
	Name     string
	Username string
	Email    string
	Password string
	Phone    *string
}

func (u *RegisterCustomerUsecase) Execute(
	ctx context.Context,
	params RegisterCustomerParams,
) (*uuid.UUID, error) {
	now := appclock.Now()

	existUsr, err := u.userRepo.
		GetByUsername(ctx, u.executor, params.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	existAcc, err := u.accountRepo.
		GetByEmail(ctx, u.executor, params.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check account: %w", err)
	}

	if existAcc != nil {
		// Ensure if the account already has password,
		// this process will not allow to overwrite it
		if existAcc.Password != "" {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category: "user_action",
				Action:   "register_customer",
				Resource: "account",
				Outcome:  applogger.OutcomeFailure,
				Metadata: map[string]any{
					"email":    params.Email,
					"username": params.Username,
					"reason":   "account already exists with password",
				},
			})
			return nil, apperrors.NewConflict(domain.ErrAccountAlreadyExists.Error())
		}

		// Check if the account already have an ID or
		// the desired username is already taken by ANOTHER user.
		if existUsr != nil &&
			existUsr.ID != existAcc.UserID {
			u.auditLogger.Log(ctx, applogger.AuditEvent{
				Category: "user_action",
				Action:   "register_customer",
				Resource: "account",
				Outcome:  applogger.OutcomeFailure,
				Metadata: map[string]any{
					"email":    params.Email,
					"username": params.Username,
					"reason":   "username already taken",
				},
			})
			return nil, apperrors.NewConflict(domain.ErrAccountAlreadyExists.Error())
		}

		hash, err := u.hasher.Hash(params.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		otpCode, err := u.otpGen.Generate()
		if err != nil {
			return nil, fmt.Errorf("failed to create otp: %w", err)
		}

		otpHash, err := u.hasher.Hash(otpCode)
		if err != nil {
			return nil, fmt.Errorf("failed to hash otp: %w", err)
		}

		challenge := domain.VerificationChallenge{
			ID:           uuid.New(),
			UserID:       &existAcc.UserID,
			Type:         domain.OTPTypeNumeric,
			Channel:      domain.OTPChannelEmail,
			Purpose:      domain.OTPPurposeRegister,
			Target:       params.Email,
			CodeHash:     otpHash,
			ExpiresAt:    now.Add(15 * time.Minute),
			AttemptCount: 0,
			CreatedAt:    now,
		}

		err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
			profileProps := userRepo.SaveProfileProps{
				UserID:    existAcc.UserID,
				Name:      &params.Name,
				Username:  &params.Username,
				Phone:     params.Phone,
				UpdatedAt: now,
			}
			if err := u.userRepo.
				SaveProfile(ctx, exec, profileProps); err != nil {
				return fmt.Errorf("failed to update user profile: %w", err)
			}

			if err := u.accountRepo.
				UpdatePasswordByUserID(
					ctx,
					exec,
					existAcc.UserID,
					hash,
				); err != nil {
				return fmt.Errorf("failed to update account password: %w", err)
			}

			if err := u.challengeRepo.
				Save(ctx, exec, challenge); err != nil {
				return fmt.Errorf("failed to save verification challenge: %w", err)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		mail := mailer.SendInput{
			To:      params.Email,
			Subject: "Verify your account",
			Text:    fmt.Sprintf("Your OTP is %s", otpCode),
		}
		if err := u.mailer.Send(mail); err != nil {
			return nil, fmt.Errorf("failed to send otp: %w", err)
		}

		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category:   "user_action",
			Action:     "register_customer",
			Resource:   "account",
			ResourceID: existAcc.ID.String(),
			Outcome:    applogger.OutcomeSuccess,
			Metadata: map[string]any{
				"email":    params.Email,
				"username": params.Username,
				"type":     "link_local",
			},
		})

		return &challenge.ID, nil
	}

	if existUsr != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "register_customer",
			Resource: "account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"email":    params.Email,
				"username": params.Username,
				"reason":   "username already exists",
			},
		})
		return nil, apperrors.NewConflict(domain.ErrAccountAlreadyExists.Error())
	}

	userProps := userRepo.CreateUserProps{
		ID:        uuid.New(),
		Name:      params.Name,
		Username:  params.Username,
		Phone:     params.Phone,
		CreatedAt: now,
	}

	hash, err := u.hasher.Hash(params.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hashed: %w", err)
	}

	acc := domain.Account{
		ID:        uuid.New(),
		UserID:    userProps.ID,
		Email:     params.Email,
		Status:    domain.AccountPending,
		Type:      domain.AccountTypeCustomer,
		Password:  hash,
		CreatedAt: now,
	}

	otp, err := u.otpGen.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to create otp: %w", err)
	}

	otpHash, err := u.hasher.Hash(otp)
	if err != nil {
		return nil, fmt.Errorf("failed to hash otp: %w", err)
	}

	challenge := domain.VerificationChallenge{
		ID:           uuid.New(),
		UserID:       &userProps.ID,
		Type:         domain.OTPTypeNumeric,
		Channel:      domain.OTPChannelEmail,
		Purpose:      domain.OTPPurposeRegister,
		Target:       params.Email,
		CodeHash:     otpHash,
		ExpiresAt:    now.Add(15 * time.Minute),
		AttemptCount: 0,
		CreatedAt:    now,
	}

	cust := customerDomain.Customer{
		ID:        uuid.New(),
		UserID:    userProps.ID,
		CreatedAt: now,
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.userRepo.
			CreateUser(ctx, exec, userProps); err != nil {
			return fmt.Errorf("failed to register: %w", err)
		}

		if err := u.customerRepo.
			Create(ctx, exec, cust); err != nil {
			return fmt.Errorf("failed to register: %w", err)
		}

		if err := u.accountRepo.
			Create(ctx, exec, acc); err != nil {
			return fmt.Errorf("failed to register: %w", err)
		}

		if err := u.challengeRepo.
			Save(ctx, exec, challenge); err != nil {
			return fmt.Errorf("failed to save challenge: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	mail := mailer.SendInput{
		To:      params.Email,
		Subject: "Verify your account",
		Text:    fmt.Sprintf("Your OTP is %s", otp),
	}
	if err := u.mailer.Send(mail); err != nil {
		return nil, fmt.Errorf("failed to send otp: %w", err)
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category:   "user_action",
		Action:     "register_customer",
		Resource:   "account",
		ResourceID: acc.ID.String(),
		Outcome:    applogger.OutcomeSuccess,
		Metadata: map[string]any{
			"email":    params.Email,
			"username": params.Username,
		},
	})

	return &challenge.ID, nil
}
