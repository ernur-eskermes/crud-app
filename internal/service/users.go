package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ernur-eskermes/crud-app/pkg/logging"

	"github.com/ernur-eskermes/crud-app/pkg/otp"
	cache "github.com/ernur-eskermes/go-homeworks/2-cache-ttl"

	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"

	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
	"github.com/ernur-eskermes/crud-app/pkg/hash"
)

type UsersRepository interface {
	Create(ctx context.Context, user *core.User) error
	GetByCredentials(ctx context.Context, email, password string) (core.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (core.User, error)
	Verify(ctx context.Context, username string) error
}

type SessionsRepository interface {
	GetByToken(ctx context.Context, token string) (core.RefreshSession, error)
	Create(ctx context.Context, session core.RefreshSession) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type UsersService struct {
	repo         UsersRepository
	sessionsRepo SessionsRepository

	auditClient AuditClient

	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	cache        cache.Cache
	otpGenerator otp.Generator

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	logger *logging.Logger

	domain string
}

func NewUsersService(repo UsersRepository, sessionsRepo SessionsRepository, auditClient AuditClient, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	accessTTL, refreshTTL time.Duration, domain string, cache cache.Cache, otpGenerator otp.Generator, logger *logging.Logger,
) *UsersService {
	return &UsersService{
		repo:            repo,
		sessionsRepo:    sessionsRepo,
		auditClient:     auditClient,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		domain:          domain,
		cache:           cache,
		logger:          logger,
		otpGenerator:    otpGenerator,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input core.AuthInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	code := s.otpGenerator.RandomSecret(6)
	s.cache.Set(input.Username, code, 10*time.Minute)
	fmt.Println(code)

	user := core.User{
		Username: input.Username,
		Password: passwordHash,
	}
	if err = s.repo.Create(ctx, &user); err != nil {
		return err
	}

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ActionRegister,
		Entity:    audit.EntityUser,
		EntityID:  user.ID.String(),
		Timestamp: time.Now(),
	}); err != nil {
		s.logger.Error("failed to send log request: ", err)
	}

	return nil
}

func (s *UsersService) Verify(ctx context.Context, username, code string) error {
	c, err := s.cache.Get(username)
	if err != nil {
		return core.ErrUserCodeIncorrect
	}

	v, ok := c.(string)
	if !ok {
		return core.ErrUserCodeIncorrect
	}

	if v != code {
		return core.ErrUserCodeIncorrect
	}

	s.cache.Delete(username)

	return s.repo.Verify(ctx, username)
}

func (s *UsersService) SignIn(ctx context.Context, input core.AuthInput) (core.Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return core.Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Username, passwordHash)
	if err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return core.Tokens{}, err
		}

		return core.Tokens{}, err
	}

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ActionLogin,
		Entity:    audit.EntityUser,
		EntityID:  user.ID.String(),
		Timestamp: time.Now(),
	}); err != nil {
		s.logger.Error("failed to send log request: ", err)
	}

	return s.createSession(ctx, user.ID)
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (core.Tokens, error) {
	session, err := s.sessionsRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, core.ErrTokenNotFound) {
			return core.Tokens{}, err
		}

		return core.Tokens{}, err
	}

	if err = s.sessionsRepo.Delete(ctx, session.UserID); err != nil {
		return core.Tokens{}, err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return core.Tokens{}, core.ErrRefreshTokenExpired
	}

	return s.createSession(ctx, session.UserID)
}

func (s *UsersService) GetByID(ctx context.Context, id uuid.UUID) (core.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return user, err
	}

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ActionGet,
		Entity:    audit.EntityUser,
		EntityID:  user.ID.String(),
		Timestamp: time.Now(),
	}); err != nil {
		s.logger.Error("failed to send log request: ", err)
	}

	return user, err
}

func (s *UsersService) createSession(ctx context.Context, userID uuid.UUID) (core.Tokens, error) {
	var (
		res core.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userID.String(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	if err = s.sessionsRepo.Create(ctx, core.RefreshSession{
		UserID:    userID,
		Token:     res.RefreshToken,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	}); err != nil {
		return core.Tokens{}, err
	}

	return res, err
}
