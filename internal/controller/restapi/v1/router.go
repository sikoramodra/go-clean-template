package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sikoramodra/go-clean-template/internal/controller/restapi/middleware"
	"github.com/sikoramodra/go-clean-template/internal/usecase"
	"github.com/sikoramodra/go-clean-template/pkg/logger"
)

// NewRoutes -.
func NewRoutes(apiV1Group chi.Router, u usecase.User, l logger.Interface) {
	v1 := &V1{u: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	// User routes
	apiV1Group.Route("/user", func(r chi.Router) {
		// // Protected routes
		r.Use(middleware.Session())
		r.Get("/profile", v1.profile)
	})
}
