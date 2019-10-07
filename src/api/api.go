package api

import (
	"context"

	"github.com/labstack/echo"
	"github.com/playneta/go-sessions/src/providers"
	"github.com/playneta/go-sessions/src/repositories"
	"github.com/playneta/go-sessions/src/services"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	// API structure represetns and handles api interaction
	API struct {
		logger         *zap.SugaredLogger
		config         *providers.Config
		userRepo       repositories.User
		accountService services.Account
		echo           *echo.Echo
	}

	// Options represetns api options
	Options struct {
		fx.In

		Logger         *zap.SugaredLogger
		Config         *providers.Config
		UserRepo       repositories.User
		AccountService services.Account
		Lc             fx.Lifecycle
	}
)

// New creates new instance of API and inject it into fx.Lifecycle
func New(opts Options) *API {
	a := &API{
		logger:         opts.Logger,
		config:         opts.Config,
		userRepo:       opts.UserRepo,
		accountService: opts.AccountService,
		echo:           echo.New(),
	}

	a.echo.HidePort = true
	a.echo.HideBanner = true

	// Endpoint
	a.echo.POST("/register", a.Register)
	a.echo.POST("/sign-in", a.SignIn)

	a.echo.GET("/profile", a.Profile, a.AuthMiddleware)

	// Start & Stop server
	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			opts.Logger.Infof("starting server at: %s", opts.Config.ListenAddr)
			go func() {
				if err := a.echo.Start(opts.Config.ListenAddr); err != nil {
					opts.Logger.Errorf("error starging server: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return a.echo.Shutdown(ctx)
		},
	})

	return a
}
