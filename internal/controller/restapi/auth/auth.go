package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/sikoramodra/go-clean-template/pkg/logger"
	"github.com/sikoramodra/go-clean-template/pkg/supertokens"
)

type userRegistrar interface {
	Register(ctx context.Context, id, email string) error
}

// Adapter - wires SuperTokens' events to the User use case.
type Adapter struct {
	users userRegistrar
	l     logger.Interface
}

// New -.
func New(users userRegistrar, l logger.Interface) *Adapter {
	return &Adapter{users: users, l: l}
}

// Hooks -.
func (a *Adapter) Hooks() supertokens.Hooks {
	return supertokens.Hooks{OnSignUp: a.onSignUp}
}

func (a *Adapter) onSignUp(ctx context.Context, id, email string) error {
	if err := a.users.Register(ctx, id, email); err != nil {
		if delErr := supertokens.DeleteUser(id); delErr != nil {
			a.l.Error(fmt.Errorf("auth - onSignUp: register+rollback failed: %w", errors.Join(err, delErr)))

			return fmt.Errorf("auth - onSignUp: register+rollback failed: %w", errors.Join(err, delErr))
		}

		a.l.Warn("auth - onSignUp: register failed, rolled back supertokens user: userID=%s err=%v", id, err)

		return fmt.Errorf("auth - onSignUp - users.Register: %w", err)
	}

	return nil
}
