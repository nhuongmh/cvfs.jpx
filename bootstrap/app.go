package bootstrap

import (
	"context"

	"github.com/nhuongmh/cfvs.jpx/pkg/database/postgresdb"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
)

type Application struct {
	Env *Env
	DB  *postgresdb.DB
}

func Init() Application {
	app := &Application{}
	app.Env = NewEnv()

	ctx := context.Background()
	db, err := postgresdb.ConnectDB(ctx, app.Env.PostgresDBUrl)

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed connect database")
	}

	app.DB = db

	err = db.Migrate()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed migrate database")
	}
	logger.Log.Info().Msg("Successfully migrated database")

	return *app
}

func (app *Application) CloseDB() {
	app.DB.Close()
}
