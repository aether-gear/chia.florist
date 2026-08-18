package usecase

import (
	"context"
	"errors"
	"time"

	applogger "service-core/internal/common/logger"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzRepo "service-core/internal/modules/authorization/repository"
	staffDomain "service-core/internal/modules/staff/domain"
	staffRepo "service-core/internal/modules/staff/repository"
	userDomain "service-core/internal/modules/user/domain"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockExecutor struct{}

func (m *mockExecutor) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("UPDATE 1"), nil
}

func (m *mockExecutor) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockExecutor) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return &mockRow{}
}

type mockRow struct{}

func (r *mockRow) Scan(dest ...any) error {
	return nil
}

type mockTransactor struct{}

func (m *mockTransactor) WithinTransaction(ctx context.Context, fn func(exec transaction.Executor) error) error {
	return fn(&mockExecutor{})
}

type mockAuditLogger struct {
	events []applogger.AuditEvent
}

func (m *mockAuditLogger) Log(ctx context.Context, event applogger.AuditEvent) {
	m.events = append(m.events, event)
}

type mockStaffRepo struct {
	staff        *staffDomain.Staff
	getByIDError error
	updateError  error
	deleteError  error
	createCalls  int
	updateCalls  int
	deleteCalls  int
}

func (m *mockStaffRepo) Create(ctx context.Context, exec transaction.Executor, staff staffDomain.Staff) error {
	m.createCalls++
	return nil
}

func (m *mockStaffRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*staffDomain.Staff, error) {
	if m.getByIDError != nil {
		return nil, m.getByIDError
	}
	return m.staff, nil
}

func (m *mockStaffRepo) GetProfileByUserID(ctx context.Context, exec transaction.Executor, userID uuid.UUID) (*staffDomain.StaffProfile, error) {
	return nil, nil
}

func (m *mockStaffRepo) FindStaff(ctx context.Context, exec transaction.Executor, params staffRepo.FindStaffParams) ([]staffDomain.StaffProfile, int, error) {
	return nil, 0, nil
}

func (m *mockStaffRepo) Update(ctx context.Context, exec transaction.Executor, staffID uuid.UUID, name string, logoUrl *string, bannerUrl *string) error {
	m.updateCalls++
	return m.updateError
}

func (m *mockStaffRepo) Delete(ctx context.Context, exec transaction.Executor, staffID uuid.UUID) error {
	m.deleteCalls++
	return m.deleteError
}

var _ staffRepo.StaffRepository = (*mockStaffRepo)(nil)

type mockStaffMembershipRepo struct {
	membership         *authzDomain.StaffMembership
	targetMembership   *authzDomain.StaffMembership
	roles              []authzDomain.Role
	accounts           []authzDomain.StaffAccountMember
	getMemError        error
	getRolesError      error
	listAccsError      error
	deleteByStaffErr   error
	deleteByAccErr     error
	saveCalls          int
	deleteByAccCalls   int
	deleteByStaffCalls int
}

func (m *mockStaffMembershipRepo) GetByAccountID(ctx context.Context, exec transaction.Executor, accountID uuid.UUID) (*authzDomain.StaffMembership, error) {
	return m.membership, nil
}

func (m *mockStaffMembershipRepo) GetByAccountIDAndStaffID(ctx context.Context, exec transaction.Executor, accountID, staffID uuid.UUID) (*authzDomain.StaffMembership, error) {
	if m.getMemError != nil {
		return nil, m.getMemError
	}
	if m.targetMembership != nil && accountID == m.targetMembership.AccountID {
		return m.targetMembership, nil
	}
	return m.membership, nil
}

func (m *mockStaffMembershipRepo) ListRolesByAccountIDAndStaffID(ctx context.Context, exec transaction.Executor, accountID, staffID uuid.UUID) ([]authzDomain.Role, error) {
	if m.getRolesError != nil {
		return nil, m.getRolesError
	}
	return m.roles, nil
}

func (m *mockStaffMembershipRepo) Save(ctx context.Context, exec transaction.Executor, membership authzDomain.StaffMembership) error {
	m.saveCalls++
	return nil
}

func (m *mockStaffMembershipRepo) ListAccountsByStaffID(ctx context.Context, exec transaction.Executor, staffID uuid.UUID) ([]authzDomain.StaffAccountMember, error) {
	if m.listAccsError != nil {
		return nil, m.listAccsError
	}
	return m.accounts, nil
}

func (m *mockStaffMembershipRepo) DeleteByAccountIDAndStaffID(ctx context.Context, exec transaction.Executor, accountID, staffID uuid.UUID) error {
	m.deleteByAccCalls++
	return m.deleteByAccErr
}

func (m *mockStaffMembershipRepo) DeleteByStaffID(ctx context.Context, exec transaction.Executor, staffID uuid.UUID) error {
	m.deleteByStaffCalls++
	return m.deleteByStaffErr
}

var _ authzRepo.StaffMembershipRepository = (*mockStaffMembershipRepo)(nil)

type mockAccountRepo struct {
	account     *authenDomain.Account
	deleteCalls int
}

func (m *mockAccountRepo) GetByEmail(ctx context.Context, exec transaction.Executor, email string) (*authenDomain.Account, error) {
	return nil, nil
}

func (m *mockAccountRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*authenDomain.Account, error) {
	return m.account, nil
}

func (m *mockAccountRepo) GetByUserID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*authenDomain.Account, error) {
	return m.account, nil
}

func (m *mockAccountRepo) ActivateByUserID(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}

func (m *mockAccountRepo) UpdatePasswordByUserID(ctx context.Context, exec transaction.Executor, id uuid.UUID, hashedPassword string) error {
	return nil
}

func (m *mockAccountRepo) Create(ctx context.Context, exec transaction.Executor, account authenDomain.Account) error {
	return nil
}

func (m *mockAccountRepo) DeleteByUserID(ctx context.Context, exec transaction.Executor, userID uuid.UUID) error {
	m.deleteCalls++
	return nil
}

func (m *mockAccountRepo) UpdateLastLoginAt(ctx context.Context, exec transaction.Executor, id uuid.UUID, lastLoginAt time.Time) error {
	return nil
}

var _ authenRepo.AccountRepository = (*mockAccountRepo)(nil)

type mockSessionRepo struct {
	revokeCalls    int
	revokedUserIDs []uuid.UUID
	err            error
}

func (m *mockSessionRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*authenDomain.Session, error) {
	return nil, nil
}

func (m *mockSessionRepo) RevokeByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}

func (m *mockSessionRepo) RevokeAllByUserID(ctx context.Context, exec transaction.Executor, userID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	m.revokeCalls++
	m.revokedUserIDs = append(m.revokedUserIDs, userID)
	return nil
}

func (m *mockSessionRepo) UpdateLastActivityByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}

func (m *mockSessionRepo) Save(ctx context.Context, exec transaction.Executor, session authenDomain.Session) error {
	return nil
}

var _ authenRepo.SessionRepository = (*mockSessionRepo)(nil)

type mockUserRepo struct {
	user        *userDomain.User
	createCalls int
}

func (m *mockUserRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*userDomain.User, error) {
	return m.user, nil
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, exec transaction.Executor, username string) (*userDomain.User, error) {
	return m.user, nil
}

func (m *mockUserRepo) CreateUser(ctx context.Context, exec transaction.Executor, props userRepo.CreateUserProps) error {
	m.createCalls++
	return nil
}

func (m *mockUserRepo) SaveProfile(ctx context.Context, exec transaction.Executor, props userRepo.SaveProfileProps) error {
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}

var _ userRepo.UserRepository = (*mockUserRepo)(nil)

type mockRoleRepo struct {
	role *authzDomain.Role
}

func (m *mockRoleRepo) GetByCode(ctx context.Context, exec transaction.Executor, code authzDomain.RoleCode) (*authzDomain.Role, error) {
	return m.role, nil
}

var _ authzRepo.RoleRepository = (*mockRoleRepo)(nil)

type mockPwHasher struct{}

func (m *mockPwHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *mockPwHasher) Compare(hash, password string) error {
	return nil
}

var _ authenRepo.PasswordHasher = (*mockPwHasher)(nil)

type mockUserDeletionService struct {
	deletedUsers []uuid.UUID
	err          error
}

func (m *mockUserDeletionService) DeleteUserRecord(ctx context.Context, exec transaction.Executor, userID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	m.deletedUsers = append(m.deletedUsers, userID)
	return nil
}

var _ authenRepo.UserDeletionService = (*mockUserDeletionService)(nil)

var errMock = errors.New("mock error")
