package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekifel/moneysaverz/internal/api"
	"github.com/ekifel/moneysaverz/internal/config"
	"github.com/ekifel/moneysaverz/internal/repository"
	"github.com/ekifel/moneysaverz/internal/server"
	"github.com/ekifel/moneysaverz/internal/service"
	"github.com/ekifel/moneysaverz/pkg/auth"
	"github.com/ekifel/moneysaverz/pkg/database/mongodb"
	"github.com/ekifel/moneysaverz/pkg/hash"
	"github.com/ekifel/moneysaverz/pkg/logger"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Errorf("error initializing config: %v", err.Error())

		return
	}

	// Dependencies
	mongoClient, err := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	if err != nil {
		logger.Errorf("error initializing mongo client: %v", err.Error())

		return
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:           repos,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL,
	})
	handlers := api.NewHandler(services, tokenManager)

	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
