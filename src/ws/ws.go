package ws

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/playneta/go-sessions/src/models"
	"github.com/playneta/go-sessions/src/repositories"
	"github.com/playneta/go-sessions/src/services"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	Websocket struct {
		logger         *zap.SugaredLogger
		config         *viper.Viper
		userRepo       repositories.User
		accountService services.Account

		hub map[string]User
		hmu sync.RWMutex
	}

	User struct {
		Model *models.User
		Conn  *websocket.Conn
	}

	Options struct {
		fx.In

		Logger         *zap.SugaredLogger
		Config         *viper.Viper
		Lc             fx.Lifecycle
		UserRepo       repositories.User
		AccountService services.Account
	}
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func New(opts Options) {
	socket := &Websocket{
		logger:         opts.Logger,
		config:         opts.Config,
		userRepo:       opts.UserRepo,
		accountService: opts.AccountService,
		hub:            make(map[string]User),
	}

	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := opts.Config.GetString("ws.addr")
			opts.Logger.Infof("starting websocket server at: %s", addr)

			go func() {
				if err := http.ListenAndServe(addr, socket); err != nil {
					opts.Logger.Errorf("error starting websocket server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func (s *Websocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Connecting
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("error creating upgrader: %v", err)
		return
	}
	defer c.Close()

	// Authentication
	token := r.URL.Query().Get("token")
	user, err := s.userRepo.FindByToken(token)
	if err != nil || user == nil {
		s.logger.Errorf("unable to find user by token '%s': %v", token, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s.hmu.Lock()
	s.hub[user.Email] = User{
		Model: user,
		Conn:  c,
	}
	s.hmu.Unlock()

	// Sending history to user
	s.logger.Info("sending history to user")
	messages, err := s.accountService.History(*user)
	if err != nil {
		s.logger.Errorf("error getting history: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, message := range messages {
		if err := c.WriteJSON(NewMessageEvent(message)); err != nil {
			s.logger.Errorf("error sending history: %v", err)
		}
	}

	// Sending all message of user join
	for _, user := range s.hub {
		if err := user.Conn.WriteJSON(Event{
			Type: "join",
			Data: MessageJoin{
				User: user.Model.Email,
			},
		}); err != nil {
			s.logger.Error("error sending join message: %v", err)
			continue
		}
	}

	// Message handling
	s.logger.Info("start listening for incoming messages")
	for {
		var msg MessageEvent
		if err := c.ReadJSON(&msg); err != nil {
			s.logger.Errorf("error reading message: %v", err)
			return
		}

		s.logger.Infof("got incoming message: %v", msg)

		// Creating message and saving it
		message, err := s.accountService.CreateMessage(*user, msg.To, msg.Text)
		if err != nil {
			s.logger.Errorf("error saving message: %v", err)
			continue
		}

		if message.Receiver != nil {
			s.hmu.RLock()
			userTo, ok := s.hub[message.Receiver.Email]
			s.hmu.RUnlock()

			if !ok {
				s.logger.Errorf("unknown user to send message: %s", msg.To)
				continue
			}

			if err := c.WriteJSON(NewMessageEvent(*message)); err != nil {
				s.logger.Errorf("error sending message: %v", err)
			}

			if err := userTo.Conn.WriteJSON(NewMessageEvent(*message)); err != nil {
				s.logger.Errorf("error sending message: %v", err)
			}
		} else {
			s.hmu.RLock()
			for _, user := range s.hub {
				if err := user.Conn.WriteJSON(NewMessageEvent(*message)); err != nil {
					s.logger.Error("error sending message: %v", err)
					continue
				}
			}
			s.hmu.RUnlock()
		}
	}
}

// func (s *Websocket) handleMessage()
