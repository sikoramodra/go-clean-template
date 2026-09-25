package restapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/metrics"
	"github.com/riandyrn/otelchi"
	"github.com/sikoramodra/go-clean-template/config"
	_ "github.com/sikoramodra/go-clean-template/docs" // Swagger docs.
	"github.com/sikoramodra/go-clean-template/internal/controller/restapi/middleware"
	v1 "github.com/sikoramodra/go-clean-template/internal/controller/restapi/v1"
	"github.com/sikoramodra/go-clean-template/internal/usecase"
	"github.com/sikoramodra/go-clean-template/pkg/logger"
	"github.com/sikoramodra/go-clean-template/pkg/supertokens"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewRouter -.
// Swagger spec:
//
//	@title       Go Clean Template API
//	@description Multi-domain clean architecture template with translation, user, and task management
//	@version     1.0
//	@host        localhost:8080
//	@BasePath    /
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
//	@tag.name auth
//	@tag.description SuperTokens-backed authentication endpoints. Served by the SuperTokens SDK middleware, not by this codebase's Go handlers.
//	@tag.name user
//	@tag.description Authenticated user profile endpoints.
func NewRouter(r chi.Router, cfg *config.Config, u usecase.User, l logger.Interface) {
	// Options
	r.Use(middleware.Logger(l))
	r.Use(middleware.Recovery(l))

	// Required by SuperTokens frontend SDK (cookie-based sessions cross-origin).
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.SuperTokens.WebsiteDomain},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   append([]string{"Content-Type"}, supertokens.AllCORSHeaders()...),
		AllowCredentials: true,
	}))

	r.Use(supertokens.Middleware)

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		r.Use(metrics.Collector(metrics.CollectorOpts{
			Host:  false,
			Proto: true,
			Skip: func(r *http.Request) bool {
				return r.Method == "OPTIONS" || r.URL.Path == "/metrics"
			},
		}))
		r.Handle("/metrics", metrics.Handler())
	}

	// Swagger
	if cfg.Swagger.Enabled {
		r.Get("/swagger/*", httpSwagger.Handler())
	}

	// K8s probe
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Routers
	r.Route("/v1", func(r chi.Router) {
		if cfg.Tracing.Enabled {
			r.Use(otelchi.Middleware("my-service-name", otelchi.WithChiRoutes(r)))
		}

		v1.NewRoutes(r, u, l)
	})
}
