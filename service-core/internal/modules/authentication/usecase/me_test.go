package usecase

import (
	"context"
	"testing"

	"service-core/internal/modules/authentication/domain"
	authorDomain "service-core/internal/modules/authorization/domain"
	userDomain "service-core/internal/modules/user/domain"
	transaction "service-core/internal/shared/transaction"
	"github.com/google/uuid"
)

type mockActorService struct {
	actor *authorDomain.Actor
}

func (m *mockActorService) Load(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ *uuid.UUID) (*authorDomain.Actor, error) {
	return m.actor, nil
}

func TestMeUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	existingAcc := &domain.Account{
		ID:     accountID,
		UserID: userID,
		Email:  "me@gmail.com",
		Status: domain.AccountActive,
	}

	avatar := "http://example.com/me.png"
	existingUser := &userDomain.User{
		ID:        userID,
		Name:      "Me User",
		AvatarURL: &avatar,
	}

	accountRepoMock := &mockAccountRepo{account: existingAcc}
	userRepoMock := &mockUserRepo{user: existingUser}
	actorSvcMock := &mockActorService{
		actor: &authorDomain.Actor{AccountID: accountID},
	}

	uc := NewMeUsecase(
		&mockExecutor{},
		accountRepoMock,
		userRepoMock,
		actorSvcMock,
	)

	authCtx := domain.AuthContext{
		UserID:          userID,
		IsAuthenticated: true,
	}

	result, err := uc.Execute(ctx, authCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected result to be returned")
	}

	if result.User == nil || result.User.AvatarURL == nil || *result.User.AvatarURL != avatar {
		t.Errorf("expected AvatarURL to be %s, got %v", avatar, result.User.AvatarURL)
	}
}
