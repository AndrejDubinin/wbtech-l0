package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndrejDubinin/wbtech-l0/internal/app"
)

func main() {
	// TODO: add logger
	ctx := runSignalHandler(context.Background())

	initOpts()
	app, err := app.NewApp(ctx, app.NewConfig(opts))
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal("{FATAL} error starting app", err)
	}
}

func runSignalHandler(ctx context.Context) context.Context {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	sigCtx, cancel := context.WithCancel(ctx)

	go func() {
		defer signal.Stop(sigterm)
		defer cancel()

		for {
			select {
			case sig, ok := <-sigterm:
				if !ok {
					log.Printf("[signal] signal chan closed: %s\n", sig.String())
					return
				}

				log.Printf("[signal] signal recv: %s\n", sig.String())
				return
			case _, ok := <-sigCtx.Done():
				if !ok {
					log.Println("[signal] context closed")
					return
				}

				log.Printf("[signal] ctx done: %s\n", ctx.Err().Error())
				return
			}
		}
	}()

	return sigCtx
}
