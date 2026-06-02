package usecase

import (
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"
	userRepo "service-core/internal/modules/user/repository"
	mailer "service-core/internal/shared/mailer"
	otp "service-core/internal/shared/otp"

	"github.com/google/uuid"
)

type RegisterUsecase struct {
	accountRepo   repository.AccountRepository
	hasher        domain.PasswordHasher
	userRepo      userRepo.UserRepository
	challengeRepo repository.VerificationChallengeRepository
	otpGen        otp.Generator
	mailer        mailer.Sender
}

func NewRegisterUsecase(
	accountRepo repository.AccountRepository,
	hasher domain.PasswordHasher,
	userRepo userRepo.UserRepository,
	challengeRepo repository.VerificationChallengeRepository,
	otpGen otp.Generator,
	mailer mailer.Sender,
) *RegisterUsecase {
	return &RegisterUsecase{
		accountRepo:   accountRepo,
		hasher:        hasher,
		userRepo:      userRepo,
		challengeRepo: challengeRepo,
		otpGen:        otpGen,
		mailer:        mailer,
	}
}

type SignUpParams struct {
	Name     string
	Username string
	Email    string
	Password string
	Phone    *string
}

func (u *RegisterUsecase) Execute(params SignUpParams) (*uuid.UUID, error) {
	now := time.Now()

	existUsr, err := u.userRepo.GetByUsername(params.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}
	existAcc, err := u.accountRepo.GetByEmail(params.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check account: %w", err)
	}

	if existAcc != nil || existUsr != nil {
		return nil, apperrors.NewConflict(domain.ErrAccountAlreadyExists.Error())
	}

	user := userRepo.CreateUserProps{
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
		UserID:    user.ID,
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
		UserID:       &user.ID,
		Type:         domain.OTPTypeNumeric,
		Channel:      domain.OTPChannelEmail,
		Purpose:      domain.OTPPurposeRegister,
		Target:       params.Email,
		CodeHash:     otpHash,
		ExpiresAt:    now.Add(15 * time.Minute),
		AttemptCount: 0,
		CreatedAt:    now,
	}

	if err := u.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	if err := u.accountRepo.Create(acc); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	if err := u.challengeRepo.Save(challenge); err != nil {
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

	return &challenge.ID, nil
}
