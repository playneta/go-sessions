package src

import (
	"github.com/playneta/go-sessions/src/api"
	"github.com/playneta/go-sessions/src/providers"
	"github.com/playneta/go-sessions/src/repositories"
	"github.com/playneta/go-sessions/src/services"
	"go.uber.org/fx"
)

// Run starting main application running fx with providers and ivoke api.New
func Run() {
	app := fx.New(
		fx.Provide(
			providers.NewConfig,
			providers.NewLogger,
			providers.NewDB,
			providers.NewBcryptHasher,
			services.NewAccount,
			repositories.NewUser,
		),

		fx.Invoke(
			api.New,
		),
	)

	app.Run()
}