package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/wbtech-l0/internal/app"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger, err := cfg.Build()
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}
	defer func() {
		log.Println("logger sync...")
		if err := logger.Sync(); err != nil {
			log.Println("logger.Sync:", err)
		}
	}()

	ctx := runSignalHandler(context.Background(), logger)

	initOpts()
	app, err := app.NewApp(ctx, app.NewConfig(opts), logger)
	if err != nil {
		logger.Fatal("{FATAL}", zap.Error(err))
	}

	if err := app.Run(ctx); err != nil {
		logger.Fatal("{FATAL} error starting app", zap.Error(err))
	}
}

func runSignalHandler(ctx context.Context, logger *zap.Logger) context.Context {
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
					logger.Info("[signal] signal chan closed", zap.String("signal", sig.String()))
					return
				}

				logger.Info("[signal] signal recv", zap.String("signal", sig.String()))
				return
			case _, ok := <-sigCtx.Done():
				if !ok {
					logger.Info("[signal] context closed")
					return
				}

				logger.Error("[signal] ctx done", zap.Error(ctx.Err()))
				return
			}
		}
	}()

	return sigCtx
}
