package usecase

import (
	"context"
	"fmt"
	"strings"

	appclock "service-core/internal/common/clock"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/staff/domain"
	"service-core/internal/modules/staff/repository"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreateStaffUsecase struct {
	staffRepo   repository.StaffRepository
	userRepo    userRepo.UserRepository
	executor    transaction.Executor
	transactor  transaction.Transactor
	auditLogger applogger.AuditLogger
}

func NewCreateStaffUsecase(
	staffRepo repository.StaffRepository,
	userRepo userRepo.UserRepository,
	executor transaction.Executor,
	transactor transaction.Transactor,
	auditLogger applogger.AuditLogger,
) *CreateStaffUsecase {
	return &CreateStaffUsecase{
		staffRepo:   staffRepo,
		userRepo:    userRepo,
		executor:    executor,
		transactor:  transactor,
		auditLogger: auditLogger,
	}
}

type CreateStaffInput struct {
	Name        string
	Description *string
	LogoUrl     *string
	BannerUrl   *string
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var sb strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			sb.WriteRune('-')
		}
	}
	return sb.String()
}

func (u *CreateStaffUsecase) Execute(
	ctx context.Context,
	input CreateStaffInput,
) error {
	now := appclock.Now()
	newUserID := uuid.New()
	newStaffID := uuid.New()

	slugName := slugify(input.Name)
	if slugName == "" {
		slugName = "staff"
	}
	username := fmt.Sprintf("%s-%s", slugName, uuid.New().String()[:8])

	newUser := userRepo.CreateUserProps{
		ID:        newUserID,
		Name:      input.Name,
		Username:  username,
		Phone:     nil,
		CreatedAt: now,
	}
	if input.LogoUrl != nil {
		newUser.AvatarURL = input.LogoUrl
	}

	staff := domain.Staff{
		ID:        newStaffID,
		UserID:    newUserID,
		CreatedAt: now,
	}

	err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.userRepo.CreateUser(ctx, exec, newUser); err != nil {
			return fmt.Errorf("failed to create user for staff: %w", err)
		}
		if err := u.staffRepo.Create(ctx, exec, staff); err != nil {
			return fmt.Errorf("failed to create staff entity: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category:   "user_action",
		Action:     "create_staff_profile",
		Resource:   "staff",
		ResourceID: staff.ID.String(),
		Outcome:    applogger.OutcomeSuccess,
		Metadata:   map[string]any{"name": input.Name, "user_id": newUserID.String()},
	})

	return nil
}
