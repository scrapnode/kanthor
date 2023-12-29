package main

import (
	_ "embed"
	"log"
	"runtime/debug"

	"github.com/scrapnode/kanthor/cmd"
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/project"
)

//go:embed .version
var version string

func main() {
	project.SetVersion(version)

	provider, err := configuration.New()
	if err != nil {
		panic(err)
	}
	command := cmd.New(provider)

	defer func() {
		if r := recover(); r != nil {
			log.Panicf("main.recover: %v | stack: %s", r, debug.Stack())
		}
	}()

	if err := command.Execute(); err != nil {
		log.Panicf("main.error: %s | stack: %s", err, debug.Stack())
	}
}
