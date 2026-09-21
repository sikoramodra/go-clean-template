// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sikoramodra/go-clean-template/config"
	"github.com/sikoramodra/go-clean-template/internal/controller/restapi"
	"github.com/sikoramodra/go-clean-template/internal/controller/restapi/auth"
	persistUserRepo "github.com/sikoramodra/go-clean-template/internal/repo/persistent/user"
	"github.com/sikoramodra/go-clean-template/internal/usecase"
	"github.com/sikoramodra/go-clean-template/internal/usecase/user"
	"github.com/sikoramodra/go-clean-template/pkg/httpserver"
	"github.com/sikoramodra/go-clean-template/pkg/logger"
	"github.com/sikoramodra/go-clean-template/pkg/postgres"
	"github.com/sikoramodra/go-clean-template/pkg/supertokens"
	"github.com/sikoramodra/go-clean-template/pkg/tracing"
)

type useCases struct {
	user usecase.User
}

type servers struct {
	http *httpserver.Server
}

func initUseCases(pg *postgres.Postgres) useCases {
	userRepo := persistUserRepo.New(pg)

	return useCases{
		user: user.New(userRepo),
	}
}

func initServers(cfg *config.Config, uc useCases, l logger.Interface) servers {
	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	restapi.NewRouter(httpServer.Router, cfg, uc.user, l)

	return servers{
		http: httpServer,
	}
}

func (s *servers) startServers() {
	s.http.Start()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	ctx := context.Background()

	// Tracing
	shutdownTracing, err := tracing.New(ctx, tracing.Config{
		Enabled:     cfg.Tracing.Enabled,
		ServiceName: cfg.App.Name,
		Version:     cfg.App.Version,
		Endpoint:    cfg.Tracing.OTLPEndpoint,
		Insecure:    cfg.Tracing.OTLPInsecure,
		SampleRate:  cfg.Tracing.SampleRate,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - tracing.New: %w", err))
	}
	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			l.Error(fmt.Errorf("app - Run - shutdownTracing: %w", err))
		}
	}()

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	uc := initUseCases(pg)

	// SuperTokens
	authAdapter := auth.New(uc.user, l)

	if err := supertokens.New(&supertokens.Config{
		ConnectionURI: cfg.SuperTokens.ConnectionURI,
		APIKey:        cfg.SuperTokens.APIKey,
		AppName:       cfg.SuperTokens.AppName,
		APIDomain:     cfg.SuperTokens.APIDomain,
		WebsiteDomain: cfg.SuperTokens.WebsiteDomain,
	}, authAdapter.Hooks()); err != nil {
		l.Fatal(fmt.Errorf("app - Run - supertokens.New: %w", err))
	}

	s := initServers(cfg, uc, l)
	s.startServers()
	s.waitForShutdown(l)
}
