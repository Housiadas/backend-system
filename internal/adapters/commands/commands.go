// Package commands contains the functionality for the set of commands
// currently supported by the CLI tooling.
package commands

import (
	"errors"
	"github.com/Housiadas/backend-system/pkg/sqldb"

	"github.com/Housiadas/backend-system/internal/config"
	"github.com/Housiadas/backend-system/pkg/logger"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("help provided")

type Config struct {
	DB      config.DB
	Version config.Version
	Auth    config.Auth
	Kafka   config.Kafka
}

type Command struct {
	DB      sqldb.Config
	Log     *logger.Logger
	Version config.Version
	Auth    config.Auth
	Kafka   config.Kafka
}

func New(
	cfg Config,
	log *logger.Logger,
	build string,
	serviceName string,
) *Command {
	return &Command{
		DB: sqldb.Config{
			User:         cfg.DB.User,
			Password:     cfg.DB.Password,
			Host:         cfg.DB.Host,
			Name:         cfg.DB.Name,
			MaxIdleConns: cfg.DB.MaxIdleConns,
			MaxOpenConns: cfg.DB.MaxOpenConns,
			DisableTLS:   cfg.DB.DisableTLS,
		},
		Log: log,
		Version: config.Version{
			Build: build,
			Desc:  serviceName,
		},
		Auth: config.Auth{
			KeysFolder: "/keys",
		},
		Kafka: cfg.Kafka,
	}
}
