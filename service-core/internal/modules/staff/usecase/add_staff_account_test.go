package usecase

import (
	"context"
	"testing"
	"time"

	apperrors "service-core/internal/common/errors"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzDomain "service-core/internal/modules/authorization/domain"
	staffDomain "service-core/internal/modules/staff/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type customMockAccountRepo struct {
	mockAccountRepo
	existingAccount *authenDomain.Account
	createCalls     int
}

func (m *customMockAccountRepo) GetByEmail(ctx context.Context, exec transaction.Executor, email string) (*authenDomain.Account, error) {
	return m.existingAccount, nil
}

func (m *customMockAccountRepo) Create(ctx context.Context, exec transaction.Executor, account authenDomain.Account) error {
	m.createCalls++
	return nil
}

func TestAddStaffAccountUsecase_Success(t *testing.T) {
	actorAccountID := uuid.New()
	actorStaffID := uuid.New()
	targetStaffID := uuid.New()
	targetUserID := uuid.New()
	roleID := uuid.New()

	memRepo := &mockStaffMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   actorStaffID,
			AccountID: actorAccountID,
		},
		roles: []authzDomain.Role{
			{ID: roleID, Code: authzDomain.RoleStaffAdmin, Name: "Staff Admin"},
		},
	}
	staffR := &mockStaffRepo{
		staff: &staffDomain.Staff{
			ID:        targetStaffID,
			UserID:    targetUserID,
			CreatedAt: time.Now(),
		},
	}
	accRepo := &customMockAccountRepo{}
	roleR := &mockRoleRepo{
		role: &authzDomain.Role{ID: roleID, Code: authzDomain.RoleStaff, Name: "Staff"},
	}
	pwHash := &mockPwHasher{}
	userR := &mockUserRepo{}
	exec := &mockExecutor{}
	tx := &mockTransactor{}
	audit := &mockAuditLogger{}

	uc := NewAddStaffAccountUsecase(exec, tx, accRepo, pwHash, userR, staffR, memRepo, roleR, audit)

	params := AddStaffAccountParams{
		ActorAccountID: actorAccountID,
		ActorStaffID:   actorStaffID,
		StaffID:        targetStaffID,
		Email:          "newstaff@chia.florist",
		Password:       "password123",
	}

	err := uc.Execute(context.Background(), params)
	require.NoError(t, err)
	assert.Equal(t, 1, accRepo.createCalls)
	assert.Equal(t, 1, memRepo.saveCalls)
	assert.Len(t, audit.events, 1)
	assert.Equal(t, "add_staff_account", audit.events[0].Action)
}

func TestAddStaffAccountUsecase_ValidationErrors(t *testing.T) {
	uc := NewAddStaffAccountUsecase(
		&mockExecutor{},
		&mockTransactor{},
		&customMockAccountRepo{},
		&mockPwHasher{},
		&mockUserRepo{},
		&mockStaffRepo{},
		&mockStaffMembershipRepo{},
		&mockRoleRepo{},
		&mockAuditLogger{},
	)

	// Missing email
	err := uc.Execute(context.Background(), AddStaffAccountParams{
		Email:    "",
		Password: "password123",
	})
	assert.Error(t, err)
	var badReq *apperrors.AppError
	assert.ErrorAs(t, err, &badReq)
	assert.Equal(t, 400, badReq.StatusCode)

	// Missing password
	err = uc.Execute(context.Background(), AddStaffAccountParams{
		Email:    "staff@chia.florist",
		Password: "",
	})
	assert.Error(t, err)
	assert.ErrorAs(t, err, &badReq)
	assert.Equal(t, 400, badReq.StatusCode)
}

func TestAddStaffAccountUsecase_NonAdminForbidden(t *testing.T) {
	actorAccountID := uuid.New()
	actorStaffID := uuid.New()
	targetStaffID := uuid.New()

	memRepo := &mockStaffMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   actorStaffID,
			AccountID: actorAccountID,
		},
		roles: []authzDomain.Role{
			{ID: uuid.New(), Code: authzDomain.RoleStaff, Name: "Staff"},
		},
	}
	staffR := &mockStaffRepo{
		staff: &staffDomain.Staff{ID: targetStaffID, UserID: uuid.New()},
	}

	uc := NewAddStaffAccountUsecase(
		&mockExecutor{},
		&mockTransactor{},
		&customMockAccountRepo{},
		&mockPwHasher{},
		&mockUserRepo{},
		staffR,
		memRepo,
		&mockRoleRepo{},
		&mockAuditLogger{},
	)

	err := uc.Execute(context.Background(), AddStaffAccountParams{
		ActorAccountID: actorAccountID,
		ActorStaffID:   actorStaffID,
		StaffID:        targetStaffID,
		Email:          "staff@chia.florist",
		Password:       "password123",
	})
	assert.Error(t, err)
	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, 403, appErr.StatusCode)
}

func TestAddStaffAccountUsecase_DuplicateEmail(t *testing.T) {
	actorAccountID := uuid.New()
	actorStaffID := uuid.New()
	targetStaffID := uuid.New()

	memRepo := &mockStaffMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   actorStaffID,
			AccountID: actorAccountID,
		},
		roles: []authzDomain.Role{
			{ID: uuid.New(), Code: authzDomain.RoleStaffAdmin, Name: "Staff Admin"},
		},
	}
	accRepo := &customMockAccountRepo{
		existingAccount: &authenDomain.Account{
			ID:    uuid.New(),
			Email: "taken@chia.florist",
		},
	}

	uc := NewAddStaffAccountUsecase(
		&mockExecutor{},
		&mockTransactor{},
		accRepo,
		&mockPwHasher{},
		&mockUserRepo{},
		&mockStaffRepo{},
		memRepo,
		&mockRoleRepo{},
		&mockAuditLogger{},
	)

	err := uc.Execute(context.Background(), AddStaffAccountParams{
		ActorAccountID: actorAccountID,
		ActorStaffID:   actorStaffID,
		StaffID:        targetStaffID,
		Email:          "taken@chia.florist",
		Password:       "password123",
	})
	assert.Error(t, err)
	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, 409, appErr.StatusCode)
}
