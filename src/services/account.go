package services

import (
	"errors"
	"strings"

	"github.com/playneta/go-sessions/src/models"
	"github.com/playneta/go-sessions/src/providers"
	"github.com/playneta/go-sessions/src/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	Account interface {
		Register(email, password string) (*models.User, error)
		Authorize(email, password string) (*models.User, error)
	}

	AccountOptions struct {
		fx.In

		Logger      *zap.SugaredLogger
		Hasher      providers.Hasher
		AccountRepo repositories.User
	}

	accountService struct {
		accountRepo repositories.User
		logger      *zap.SugaredLogger
		hasher      providers.Hasher
	}
)

var (
	ErrMalformedEmail  = errors.New("malformed email")
	ErrPasswordToSmall = errors.New("pussword must be more than 5 symbols")
	ErrUnauthorized    = errors.New("unauthorized")
)

func NewAccount(opts AccountOptions) Account {
	return &accountService{
		logger:      opts.Logger.Named("account_service"),
		accountRepo: opts.AccountRepo,
		hasher:      opts.Hasher,
	}
}

func (a *accountService) Register(email, password string) (*models.User, error) {
	if !strings.Contains(email, "@") {
		return nil, ErrMalformedEmail
	}

	if len(password) < 6 {
		return nil, ErrPasswordToSmall
	}

	hashedPassword, err := a.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	return a.accountRepo.Create(email, hashedPassword)
}

func (a *accountService) Authorize(email, password string) (*models.User, error) {
	user, err := a.accountRepo.FindByEmail(email)
	if err != nil || user == nil {
		a.logger.Debugf("user not found")
		return nil, ErrUnauthorized
	}

	if !a.hasher.Compare(password, user.Password) {
		a.logger.Debugf("password hash not match")
		return nil, ErrUnauthorized
	}

	if err := a.accountRepo.UpdateToken(user, generateToken()); err != nil {
		a.logger.Errorf("error updating user token: %v", err)
		return nil, ErrUnauthorized
	}

	return user, nil
}

func generateToken() string {
	return "secure_token"
}
