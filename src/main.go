package src

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/playneta/go-sessions/src/api"
	"github.com/playneta/go-sessions/src/providers"
	"github.com/playneta/go-sessions/src/repositories"
	"github.com/playneta/go-sessions/src/services"
	"github.com/playneta/go-sessions/src/ws"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
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
			repositories.NewMessage,
		),

		fx.Invoke(
			api.New,
			ws.New,
		),
	)

	app.Run()
}

func Migrate(dir string) {
	app := fx.New(
		fx.Provide(
			providers.NewConfig,
			providers.NewLogger,
		),
		fx.Invoke(func(v *viper.Viper, logger *zap.SugaredLogger) {
			if err := migrate(dir, v); err != nil {
				logger.Errorf("error migrating: %v", err)
			}

			logger.Infof("migration successfull")
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func migrate(dir string, v *viper.Viper) error {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		v.GetString("db.user"),
		v.GetString("db.password"),
		v.GetString("db.addr"),
		v.GetString("db.db"),
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}

	if err := goose.Run("up", db, dir); err != nil {
		return err
	}

	return nil
}
