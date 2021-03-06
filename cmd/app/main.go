package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ernur-eskermes/crud-app/pkg/otp"

	_ "github.com/ernur-eskermes/crud-app/docs"
	"github.com/ernur-eskermes/crud-app/internal/config"
	"github.com/ernur-eskermes/crud-app/internal/repository"
	"github.com/ernur-eskermes/crud-app/internal/service"
	"github.com/ernur-eskermes/crud-app/internal/transport/rest"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
	"github.com/ernur-eskermes/crud-app/pkg/database/postgresql"
	"github.com/ernur-eskermes/crud-app/pkg/hash"
	"github.com/ernur-eskermes/crud-app/pkg/logging"
	cache "github.com/ernur-eskermes/go-homeworks/2-cache-ttl"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	_ "github.com/lib/pq"
)

const configsDir = "configs"

// @title CRUD API
// @version 1.0
// @description REST API for CRUD App

// @host localhost:8000
// @BasePath /api/v1/

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

// Run initializes whole application.
func main() {
	logger := logging.GetLogger()

	cfg, err := config.Init(configsDir)
	if err != nil {
		logger.Fatal(err)
	}

	// Dependencies
	memCache := cache.New()

	otpGenerator := otp.NewGOTPGenerator()

	hasher := hash.NewSHA256Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Fatal(err)
	}

	validation := validator.New()

	// init db
	db, err := postgresql.NewClient(context.TODO(), 5, postgresql.StorageConfig{
		ConnStr: cfg.Postgres.ConnStr,
		Logger:  logrusadapter.NewLogger(logger),
	})
	if err != nil {
		logger.Fatal(err)
	}

	// init deps

	repos := repository.NewRepositories(db)
	services := service.NewServices(service.Deps{
		Repos:          repos,
		Hasher:         hasher,
		Cache:          memCache,
		OtpGenerator:   otpGenerator,
		TokenManager:   tokenManager,
		AccessTokenTTL: cfg.Auth.JWT.AccessTokenTTL,
		Environment:    cfg.Environment,
		Domain:         cfg.HTTP.Host,
	})
	handlers := rest.NewHandler(services, tokenManager, validation, logger)

	// init & run server

	app := fiber.New(fiber.Config{
		WriteTimeout: cfg.HTTP.WriteTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		BodyLimit:    cfg.HTTP.MaxHeaderMegabytes << 20,
	})

	handlers.InitRouter(app, cfg)

	go func() {
		if err := app.Listen(":" + cfg.HTTP.Port); err != nil {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("Shutting down server")

	if err := app.Shutdown(); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	db.Close()
}
