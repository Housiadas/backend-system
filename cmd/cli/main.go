package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"

	"github.com/Housiadas/backend-system/internal/app/commands"
	"github.com/Housiadas/backend-system/pkg/logger"
)

var build = "develop"

func main() {
	if err := run(); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			fmt.Println("msg", err)
		}
		os.Exit(1)
	}
}

func run() error {
	// -------------------------------------------------------------------------
	// Initialize Configuration
	// -------------------------------------------------------------------------
	c, err := LoadConfig("../../")
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// Initialize Logger
	// -------------------------------------------------------------------------
	traceIDFn := func(context.Context) string {
		return "00000000-0000-0000-0000-000000000000"
	}
	requestIDFn := func(context.Context) string {
		return "00000000-0000-0000-0000-000000000000"
	}
	log := logger.New(io.Discard, logger.LevelInfo, "CMD", traceIDFn, requestIDFn)

	// -------------------------------------------------------------------------
	// Initialize commands
	// -------------------------------------------------------------------------
	cmd := commands.New(c, log, build, "CMD")

	return processCommands(os.Args, cmd)
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (cfg commands.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile("config.yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}

// processCommands handles the execution of the commands specified on the command line.
func processCommands(args []string, cmd *commands.Command) error {
	switch args[1] {
	case "useradd":
		name := args[2]
		email := args[3]
		password := args[4]
		if err := cmd.UserAdd(name, email, password); err != nil {
			return fmt.Errorf("adding user: %w", err)
		}

	case "genkey":
		if err := cmd.GenKey(); err != nil {
			return fmt.Errorf("key generation: %w", err)
		}

	default:
		fmt.Println("seed:       add data to the database")
		fmt.Println("useradd:    add a new user to the database")
		fmt.Println("genkey:     generate a set of private/public key files")
		fmt.Println("gentoken:	 generate a JWT for a user with claims")
		fmt.Println("userevents: kafka consumer to listen to user events")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}
