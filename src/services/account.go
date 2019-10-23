package services

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/m1ome/randstr"
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
		CreateMessage(user models.User, receiverEmail, text string) (*models.Message, error)
		History(user models.User) ([]models.Message, error)
	}

	AccountOptions struct {
		fx.In

		Logger      *zap.SugaredLogger
		Hasher      providers.Hasher
		AccountRepo repositories.User
		MessageRepo repositories.Message
	}

	accountService struct {
		accountRepo repositories.User
		messageRepo repositories.Message
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
		messageRepo: opts.MessageRepo,
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

func (a *accountService) CreateMessage(user models.User, receiverEmail, text string) (*models.Message, error) {
	if len(text) == 0 {
		return nil, errors.New("empty message text")
	}

	var receiverID int64
	if receiverEmail != "" {
		r, err := a.accountRepo.FindByEmail(receiverEmail)
		if err != nil {
			return nil, err
		}

		receiverID = r.ID
	}

	message := &models.Message{
		UserId:     user.ID,
		ReceiverId: receiverID,
		Text:       text,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := a.messageRepo.Create(message); err != nil {
		return nil, err
	}

	return message, nil
}

func (a *accountService) History(user models.User) ([]models.Message, error) {
	public, err := a.messageRepo.LastPublicMessages(10)
	if err != nil {
		return nil, err
	}

	private, err := a.messageRepo.LastPrivateMessages(user, 10)
	if err != nil {
		return nil, err
	}

	messages := make([]models.Message, 0)
	messages = append(messages, public...)
	messages = append(messages, private...)

	sort.Slice(messages, func(i, j int) bool {
		return messages[j].CreatedAt.After(messages[i].CreatedAt)
	})

	return messages, nil
}

func generateToken() string {
	return randstr.GetString(16)
}
