package main

import (
	"github.com/scrapnode/kanthor/cmd"
	"log"
	"os"
	"runtime/debug"
)

func main() {
	command := cmd.New()

	defer func() {
		if r := recover(); r != nil {
			log.Println("main.recover:", r)
			log.Println("main.recover.trace:", string(debug.Stack()))
		}
	}()

	if err := command.Execute(); err != nil {
		log.Println("main.error:", err.Error())
		os.Exit(1)
	}
}