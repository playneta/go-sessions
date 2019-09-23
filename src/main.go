package src

import (
	"github.com/playneta/go-sessions/src/api"
	"github.com/playneta/go-sessions/src/providers"
	"github.com/playneta/go-sessions/src/repositories"
	"github.com/playneta/go-sessions/src/services"
	"go.uber.org/fx"
)

func Run() {
	app := fx.New(
		fx.Provide(
			providers.NewConfig,
			providers.NewLogger,
			providers.NewDB,
			services.NewAccount,
			repositories.NewUser,
		),

		fx.Invoke(
			api.New,
		),
	)

	app.Run()
}
