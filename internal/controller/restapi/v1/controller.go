package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/sikoramodra/go-clean-template/internal/usecase"
	"github.com/sikoramodra/go-clean-template/pkg/logger"
)

// V1 -.
type V1 struct {
	u usecase.User
	l logger.Interface
	v *validator.Validate
}
