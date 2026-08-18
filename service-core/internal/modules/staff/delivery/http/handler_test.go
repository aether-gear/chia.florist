package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	applogger "service-core/internal/common/logger"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	staffDomain "service-core/internal/modules/staff/domain"
	staffRepo "service-core/internal/modules/staff/repository"
	"service-core/internal/modules/staff/usecase"
	userDomain "service-core/internal/modules/user/domain"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockExec struct{}

func (m *mockExec) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag("UPDATE 1"), nil
}
func (m *mockExec) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *mockExec) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return &mockR{}
}

type mockR struct{}

func (r *mockR) Scan(dest ...any) error { return nil }

type mockTx struct{}

func (m *mockTx) WithinTransaction(ctx context.Context, fn func(exec transaction.Executor) error) error {
	return fn(&mockExec{})
}

type mockAuditor struct{}

func (m *mockAuditor) Log(ctx context.Context, event applogger.AuditEvent) {}

type testStaffRepo struct {
	staff *staffDomain.Staff
}

func (r *testStaffRepo) Create(ctx context.Context, exec transaction.Executor, staff staffDomain.Staff) error {
	return nil
}
func (r *testStaffRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*staffDomain.Staff, error) {
	return r.staff, nil
}
func (r *testStaffRepo) GetProfileByUserID(ctx context.Context, exec transaction.Executor, userID uuid.UUID) (*staffDomain.StaffProfile, error) {
	return nil, nil
}
func (r *testStaffRepo) FindStaff(ctx context.Context, exec transaction.Executor, params staffRepo.FindStaffParams) ([]staffDomain.StaffProfile, int, error) {
	return nil, 0, nil
}
func (r *testStaffRepo) Update(ctx context.Context, exec transaction.Executor, staffID uuid.UUID, name string, logoUrl *string, bannerUrl *string) error {
	return nil
}
func (r *testStaffRepo) Delete(ctx context.Context, exec transaction.Executor, staffID uuid.UUID) error {
	return nil
}

type testMembershipRepo struct {
	membership *authzDomain.StaffMembership
	roles      []authzDomain.Role
	accounts   []authzDomain.StaffAccountMember
}

func (r *testMembershipRepo) GetByAccountID(ctx context.Context, exec transaction.Executor, accountID uuid.UUID) (*authzDomain.StaffMembership, error) {
	return r.membership, nil
}
func (r *testMembershipRepo) GetByAccountIDAndStaffID(ctx context.Context, exec transaction.Executor, accountID, staffID uuid.UUID) (*authzDomain.StaffMembership, error) {
	return r.membership, nil
}
func (r *testMembershipRepo) ListRolesByAccountIDAndStaffID(ctx context.Context, exec transaction.Executor, accountID, staffID uuid.UUID) ([]authzDomain.Role, error) {
	return r.roles, nil
}
func (r *testMembershipRepo) Save(ctx context.Context, exec transaction.Executor, membership authzDomain.StaffMembership) error {
	return nil
}
func (r *testMembershipRepo) ListAccountsByStaffID(ctx context.Context, exec transaction.Executor, staffID uuid.UUID) ([]authzDomain.StaffAccountMember, error) {
	return r.accounts, nil
}
func (r *testMembershipRepo) DeleteByAccountIDAndStaffID(ctx context.Context, exec transaction.Executor, accountID, staffID uuid.UUID) error {
	return nil
}
func (r *testMembershipRepo) DeleteByStaffID(ctx context.Context, exec transaction.Executor, staffID uuid.UUID) error {
	return nil
}

type testAccountRepo struct {
	account       *authenDomain.Account
	accountByUser *authenDomain.Account
}

func (r *testAccountRepo) GetByEmail(ctx context.Context, exec transaction.Executor, email string) (*authenDomain.Account, error) {
	return nil, nil
}
func (r *testAccountRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*authenDomain.Account, error) {
	return r.account, nil
}
func (r *testAccountRepo) GetByUserID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*authenDomain.Account, error) {
	return r.accountByUser, nil
}
func (r *testAccountRepo) ActivateByUserID(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}
func (r *testAccountRepo) UpdatePasswordByUserID(ctx context.Context, exec transaction.Executor, id uuid.UUID, hashedPassword string) error {
	return nil
}
func (r *testAccountRepo) Create(ctx context.Context, exec transaction.Executor, account authenDomain.Account) error {
	return nil
}
func (r *testAccountRepo) DeleteByUserID(ctx context.Context, exec transaction.Executor, userID uuid.UUID) error {
	return nil
}
func (r *testAccountRepo) UpdateLastLoginAt(ctx context.Context, exec transaction.Executor, id uuid.UUID, lastLoginAt time.Time) error {
	return nil
}

type testUserRepo struct {
	user *userDomain.User
}

func (r *testUserRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*userDomain.User, error) {
	return r.user, nil
}
func (r *testUserRepo) GetByUsername(ctx context.Context, exec transaction.Executor, username string) (*userDomain.User, error) {
	return nil, nil
}
func (r *testUserRepo) CreateUser(ctx context.Context, exec transaction.Executor, props userRepo.CreateUserProps) error {
	return nil
}
func (r *testUserRepo) SaveProfile(ctx context.Context, exec transaction.Executor, props userRepo.SaveProfileProps) error {
	return nil
}
func (r *testUserRepo) Delete(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}

type testUserDeletionService struct{}

func (s *testUserDeletionService) DeleteUserRecord(ctx context.Context, exec transaction.Executor, userID uuid.UUID) error {
	return nil
}

type testSessionRepo struct{}

func (r *testSessionRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*authenDomain.Session, error) {
	return nil, nil
}
func (r *testSessionRepo) RevokeByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}
func (r *testSessionRepo) RevokeAllByUserID(ctx context.Context, exec transaction.Executor, userID uuid.UUID) error {
	return nil
}
func (r *testSessionRepo) UpdateLastActivityByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) error {
	return nil
}
func (r *testSessionRepo) Save(ctx context.Context, exec transaction.Executor, session authenDomain.Session) error {
	return nil
}

type testRoleRepo struct {
	role *authzDomain.Role
}

func (r *testRoleRepo) GetByCode(ctx context.Context, exec transaction.Executor, code authzDomain.RoleCode) (*authzDomain.Role, error) {
	return r.role, nil
}

type testHasher struct{}

func (h *testHasher) Hash(p string) (string, error) { return "hash", nil }
func (h *testHasher) Compare(hash, p string) error  { return nil }

func setupTestHandler(staffID, accountID uuid.UUID) (*staffHandler, *testStaffRepo, *testMembershipRepo) {
	sRepo := &testStaffRepo{
		staff: &staffDomain.Staff{
			ID:     staffID,
			UserID: uuid.New(),
		},
	}
	mRepo := &testMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   staffID,
			AccountID: accountID,
		},
		roles: []authzDomain.Role{
			{ID: uuid.New(), Code: authzDomain.RoleStaffAdmin, Name: "Staff Admin"},
		},
		accounts: []authzDomain.StaffAccountMember{
			{
				AccountID: accountID,
				UserID:    uuid.New(),
				Email:     "jane@chia.florist",
				Name:      "Jane Doe",
				Username:  "janedoe",
				Role: authzDomain.Role{
					ID:   uuid.New(),
					Code: authzDomain.RoleStaffAdmin,
					Name: "Staff Admin",
				},
				CreatedAt: time.Now(),
			},
		},
	}
	aRepo := &testAccountRepo{
		account: &authenDomain.Account{
			ID:     uuid.New(),
			UserID: uuid.New(),
		},
	}
	uRepo := &testUserRepo{}
	rRepo := &testRoleRepo{
		role: &authzDomain.Role{ID: uuid.New(), Code: authzDomain.RoleStaff, Name: "Staff"},
	}
	hasher := &testHasher{}
	sessionRepo := &testSessionRepo{}

	exec := &mockExec{}
	tx := &mockTx{}
	audit := &mockAuditor{}
	userDeletionSvc := &testUserDeletionService{}

	createUC := usecase.NewCreateStaffUsecase(sRepo, uRepo, exec, tx, audit)
	addUC := usecase.NewAddStaffAccountUsecase(exec, tx, aRepo, hasher, uRepo, sRepo, mRepo, rRepo, audit)
	listUC := usecase.NewListStaffAccountsUsecase(exec, sRepo, mRepo, audit)
	updateUC := usecase.NewUpdateStaffUsecase(exec, tx, sRepo, mRepo, audit)
	deleteUC := usecase.NewDeleteStaffUsecase(exec, tx, sRepo, mRepo, userDeletionSvc, audit)
	removeUC := usecase.NewRemoveStaffAccountUsecase(exec, tx, sRepo, mRepo, aRepo, sessionRepo, audit)

	handler := NewStaffHandler(addUC, createUC, nil, listUC, updateUC, deleteUC, removeUC)
	return handler, sRepo, mRepo
}

func withActorContext(r *http.Request, accountID, staffID uuid.UUID) *http.Request {
	actor := &authzDomain.Actor{
		AccountID: accountID,
		StaffID:   &staffID,
		Roles: []authzDomain.Role{
			{Code: authzDomain.RoleStaffAdmin},
		},
	}
	return r.WithContext(authzSvc.WithActor(r.Context(), actor))
}

func TestHandler_ListStaffAccounts(t *testing.T) {
	staffID := uuid.New()
	accountID := uuid.New()
	handler, _, _ := setupTestHandler(staffID, accountID)

	req := httptest.NewRequest(http.MethodGet, "/staff/"+staffID.String()+"/accounts", nil)
	req = withActorContext(req, accountID, staffID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("staffID", staffID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	err := handler.ListStaffAccounts(rec, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rec.Code)
	}

	var resp listStaffAccountsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Total != 1 || resp.StaffID != staffID {
		t.Errorf("unexpected response content: %+v", resp)
	}
}

func TestHandler_UpdateStaff(t *testing.T) {
	staffID := uuid.New()
	accountID := uuid.New()
	handler, _, _ := setupTestHandler(staffID, accountID)

	body, _ := json.Marshal(updateStaffRequest{
		Name: "New Branch Name",
	})
	req := httptest.NewRequest(http.MethodPut, "/staff/"+staffID.String(), bytes.NewReader(body))
	req = withActorContext(req, accountID, staffID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("staffID", staffID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	err := handler.UpdateStaff(rec, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rec.Code)
	}
}

func TestHandler_DeleteStaff(t *testing.T) {
	staffID := uuid.New()
	accountID := uuid.New()
	handler, _, _ := setupTestHandler(staffID, accountID)

	req := httptest.NewRequest(http.MethodDelete, "/staff/"+staffID.String(), nil)
	req = withActorContext(req, accountID, staffID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("staffID", staffID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	err := handler.DeleteStaff(rec, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rec.Code)
	}
}

func TestHandler_RemoveStaffAccount(t *testing.T) {
	staffID := uuid.New()
	actorAccountID := uuid.New()
	targetAccountID := uuid.New()
	handler, _, _ := setupTestHandler(staffID, actorAccountID)

	req := httptest.NewRequest(http.MethodDelete, "/staff/"+staffID.String()+"/accounts/"+targetAccountID.String(), nil)
	req = withActorContext(req, actorAccountID, staffID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("staffID", staffID.String())
	rctx.URLParams.Add("accountID", targetAccountID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	err := handler.RemoveStaffAccount(rec, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rec.Code)
	}
}

func TestHandler_CreateStaff(t *testing.T) {
	staffID := uuid.New()
	accountID := uuid.New()
	handler, _, _ := setupTestHandler(staffID, accountID)

	body, _ := json.Marshal(createStaffRequest{
		Name:     "Floral Logistics",
		Username: "floral-logistics",
	})
	req := httptest.NewRequest(http.MethodPost, "/staff", bytes.NewReader(body))
	req = withActorContext(req, accountID, staffID)

	rec := httptest.NewRecorder()
	err := handler.CreateStaff(rec, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got: %d", rec.Code)
	}
}

func TestHandler_AddStaffAccount(t *testing.T) {
	staffID := uuid.New()
	accountID := uuid.New()
	handler, _, _ := setupTestHandler(staffID, accountID)

	body, _ := json.Marshal(addStaffAccountRequest{
		Email:    "new@chia.florist",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/staff/"+staffID.String()+"/accounts", bytes.NewReader(body))
	req = withActorContext(req, accountID, staffID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("staffID", staffID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	err := handler.AddStaffAccount(rec, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got: %d", rec.Code)
	}
}
