package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/playneta/go-sessions/src/repositories"
	"github.com/playneta/go-sessions/src/services"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	// API structure represetns and handles api interaction
	API struct {
		logger         *zap.SugaredLogger
		config         *viper.Viper
		userRepo       repositories.User
		accountService services.Account
		echo           *echo.Echo
	}

	// Options represetns api options
	Options struct {
		fx.In

		Logger         *zap.SugaredLogger
		Config         *viper.Viper
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

	// CORS
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	// Endpoint
	a.echo.POST("/register", a.Register)
	a.echo.POST("/sign-in", a.SignIn)

	a.echo.GET("/profile", a.Profile, a.AuthMiddleware)

	// Start & Stop server
	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := opts.Config.GetString("api.addr")
			opts.Logger.Infof("starting server at: %s", addr)
			go func() {
				if err := a.echo.Start(addr); err != nil {
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
