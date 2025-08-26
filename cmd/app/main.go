package main

import (
	"log"

	"github.com/AndrejDubinin/wbtech-l0/internal/app"
)

func main() {
	initOpts()
	app, err := app.NewApp(app.NewConfig(opts))
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}
	defer app.Close()

	err = app.Run()
	if err != nil {
		log.Fatalf("error starting app: %s\n", err)
	}
}
