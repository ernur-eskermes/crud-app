package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	rabbitmqClient "github.com/ernur-eskermes/crud-app/internal/transport/rabbitmq"

	"github.com/ernur-eskermes/crud-app/internal/storage/psql"

	"github.com/ernur-eskermes/crud-app/pkg/otp"

	_ "github.com/ernur-eskermes/crud-app/docs"
	"github.com/ernur-eskermes/crud-app/internal/config"
	"github.com/ernur-eskermes/crud-app/internal/service"
	"github.com/ernur-eskermes/crud-app/internal/transport/rest"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
	"github.com/ernur-eskermes/crud-app/pkg/database/postgresql"
	"github.com/ernur-eskermes/crud-app/pkg/hash"
	"github.com/ernur-eskermes/crud-app/pkg/logging"
	cache "github.com/ernur-eskermes/go-homeworks/2-cache-ttl"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	_ "github.com/lib/pq"
)

const configsDir = "configs"

// @title CRUD API
// @version 1.0
// @description REST API for CRUD App

// @host localhost:8000
// @BasePath /api/

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

	// init db
	db, err := postgresql.NewClient(context.TODO(), 5, postgresql.StorageConfig{
		ConnStr: cfg.Postgres.ConnStr,
		Logger:  logrusadapter.NewLogger(logger),
	})
	if err != nil {
		logger.Fatal(err)
	}

	// init deps

	// auditClient, err := grpcClient.NewClient(cfg.GRPC.AuditURL)

	auditClient, err := rabbitmqClient.NewClient(cfg.AMQP.URI)
	if err != nil {
		logger.Fatal(err)
	}

	usersRepo := psql.NewUsersRepo(db)
	sessionsRepo := psql.NewSessionsRepo(db)
	usersService := service.NewUsersService(usersRepo, sessionsRepo, auditClient, hasher, tokenManager, cfg.Auth.JWT.AccessTokenTTL, cfg.Auth.JWT.RefreshTokenTTL, cfg.HTTP.Host, memCache, otpGenerator, logger)

	booksRepo := psql.NewBooksRepo(db)
	booksService := service.NewBooksService(booksRepo, auditClient, tokenManager, logger)

	handlers := rest.NewHandler(usersService, booksService, tokenManager, logger)

	// init & run server

	app := fiber.New(fiber.Config{
		WriteTimeout: cfg.HTTP.WriteTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		BodyLimit:    cfg.HTTP.MaxHeaderMegabytes << 20,
	})

	handlers.InitRouter(app)

	go func() {
		if err = app.Listen(":" + cfg.HTTP.Port); err != nil {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logger.Info("Shutting down server")

	if err = app.Shutdown(); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err = auditClient.CloseConnection(); err != nil {
		logger.Errorf("failed to close audit client connection: %v", err)
	}

	db.Close()
}
