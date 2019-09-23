package services

import (
	"errors"
	"strings"

	"github.com/playneta/go-sessions/src/models"
	"github.com/playneta/go-sessions/src/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type (
	Account interface {
		Register(email, password string) (*models.User, error)
		Authorize(email, password string) (*models.User, error)
	}

	AccountOptions struct {
		fx.In

		Logger      *zap.SugaredLogger
		AccountRepo repositories.User
	}

	accountService struct {
		accountRepo repositories.User
		logger      *zap.SugaredLogger
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
	}
}

func (a *accountService) Register(email, password string) (*models.User, error) {
	if !strings.Contains(email, "@") {
		return nil, ErrMalformedEmail
	}

	if len(password) < 6 {
		return nil, ErrPasswordToSmall
	}

	hashedPassword, err := hashPassword(password)
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

	if !checkPasswordHash(password, user.Password) {
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

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
