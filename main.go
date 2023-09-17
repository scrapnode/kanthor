package main

import (
	_ "embed"
	"os"
	"runtime/debug"

	"github.com/scrapnode/kanthor/cmd"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

//go:embed .version
var version string

func main() {
	provider, err := configuration.New()
	if err != nil {
		panic(err)
	}
	conf, err := config.New(provider)
	if err != nil {
		panic(err)
	}
	conf.Version = version

	logger, err := logging.New(&conf.Logger)
	if err != nil {
		panic(err)
	}

	command := cmd.New(provider, conf, logger)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("main.recover: %v", r)
			logger.Errorf("main.recover.trace: %s", debug.Stack())
		}
	}()

	if err := command.Execute(); err != nil {
		logger.Errorf("main.error: %v", err)
		os.Exit(1)
	}
}
