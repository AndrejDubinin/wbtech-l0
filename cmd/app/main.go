package main

import (
	"log"

	"github.com/AndrejDubinin/wbtech-l0/internal/app"
)

func main() {
	// TODO: add logger
	// TODO: add graceful shutdown
	// TODO: add one context
	initOpts()
	app, err := app.NewApp(app.NewConfig(opts))
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal("{FATAL} error starting app", err)
	}
}
