package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ernur-eskermes/crud-app/pkg/otp"
	cache "github.com/ernur-eskermes/go-homeworks/2-cache-ttl"

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

type UsersService struct {
	repo         UsersRepository
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	cache        cache.Cache
	otpGenerator otp.Generator

	accessTokenTTL time.Duration

	domain string
}

func NewUsersService(repo UsersRepository, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	accessTTL time.Duration, domain string, cache cache.Cache, otpGenerator otp.Generator,
) *UsersService {
	return &UsersService{
		repo:           repo,
		hasher:         hasher,
		tokenManager:   tokenManager,
		accessTokenTTL: accessTTL,
		domain:         domain,
		cache:          cache,
		otpGenerator:   otpGenerator,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input UserSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	code := s.otpGenerator.RandomSecret(6)
	s.cache.Set(input.Username, code, 10*time.Minute)
	fmt.Println(code)

	return s.repo.Create(ctx, &core.User{
		Username: input.Username,
		Password: passwordHash,
	})
}

func (s *UsersService) Verify(ctx context.Context, username, code string) error {
	c, err := s.cache.Get(username)
	if err != nil {
		return core.ErrUserCodeExpired
	}

	v, ok := c.(string)
	if !ok {
		return core.ErrUserCodeUnknownType
	}

	if v != code {
		return core.ErrUserCodeIncorrect
	}

	s.cache.Delete(username)

	return s.repo.Verify(ctx, username)
}

func (s *UsersService) SignIn(ctx context.Context, input UserSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Username, passwordHash)
	if err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return Tokens{}, err
		}

		return Tokens{}, err
	}

	return s.createSession(user.ID.String())
}

func (s *UsersService) GetByID(ctx context.Context, id uuid.UUID) (core.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UsersService) createSession(userID string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userID, s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	return res, err
}
