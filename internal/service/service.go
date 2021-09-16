package service

import (
	"context"
	"time"

	"github.com/ekifel/moneysaverz/internal/repository"
	"github.com/ekifel/moneysaverz/pkg/auth"
	"github.com/ekifel/moneysaverz/pkg/hash"
)

type Services struct {
	Users Users
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		Users: usersService,
	}
}

type Deps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type UserSignUpInput struct {
	Username string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
