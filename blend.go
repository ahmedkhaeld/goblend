package goblend

import (
	"os"

	"github.com/ahmedkhaeld/goblend/HTTP"
	"github.com/ahmedkhaeld/goblend/config"
	"github.com/ahmedkhaeld/goblend/db"
	"github.com/ahmedkhaeld/goblend/paths"
	"github.com/rs/zerolog"
)

var (
	RequestPerSecond = 100
)

type Blender struct {
	AppName  string
	Debug    bool
	Version  string
	RootPath string
	Log      zerolog.Logger
	Server   *HTTP.Server
	DB       *db.Database
}

func (b *Blender) Blend(rootPath string, envfile string) error {

	config.Load(rootPath, envfile)

	err := paths.InitAppPaths(rootPath)
	if err != nil {
		return err
	}
	b.Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	b.Log.Log().Msg("It's Go Time!")

	b.DB = db.NewDatabase()
	err = b.DB.Connect()
	if err != nil {
		b.Log.Error().Err(err).Msg("Failed to connect to database")
		return err
	}
	defer b.DB.Close()

	b.Server = HTTP.NewServer(RequestPerSecond)

	return nil
}

func (b *Blender) Close() {
	if b.DB != nil {
		b.DB.Close()
	}
}
