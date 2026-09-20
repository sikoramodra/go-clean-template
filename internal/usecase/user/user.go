package user

import (
	"context"
	"fmt"
	"time"

	"github.com/sikoramodra/go-clean-template/internal/entity"
	"github.com/sikoramodra/go-clean-template/internal/repo"
	"github.com/sikoramodra/go-clean-template/internal/usecase"
)

// UseCase -.
type UseCase struct {
	repo repo.UserRepo
}

// New returns a User usecase instrumented with OpenTelemetry tracing spans.
func New(r repo.UserRepo) usecase.User {
	return newTraced(&UseCase{
		repo: r,
	})
}

// Register -.
func (uc *UseCase) Register(ctx context.Context, id, email string) error {
	now := time.Now().UTC()

	user := entity.User{
		ID:        id,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := uc.repo.Store(ctx, &user)
	if err != nil {
		return fmt.Errorf("UserUseCase - Register - uc.repo.Store: %w", err)
	}

	return nil
}

// GetUser -.
func (uc *UseCase) GetUser(ctx context.Context, userID string) (entity.User, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - GetUser - uc.repo.GetByID: %w", err)
	}

	return user, nil
}
