package app

import (
	"context"
	"log"
	"sync"

	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka/consumer"
	"github.com/IBM/sarama"
)

type (
	cons interface {
		ConsumeTopic(ctx context.Context, handler func(*sarama.ConsumerMessage), wg *sync.WaitGroup) error
		Close() error
	}

	App struct {
		config   config
		consumer cons
	}
)

func NewApp(config config) (*App, error) {
	cons, err := consumer.NewConsumer(config.kafka, config.consumer,
		consumer.WithReturnErrorsEnabled(true),
	)
	if err != nil {
		return nil, err
	}

	return &App{
		config:   config,
		consumer: cons,
	}, nil
}

func (a *App) Run() error {
	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
	)

	log.Printf("consumer reads topic: %s\n", a.config.consumer.Topic)
	err := a.consumer.ConsumeTopic(ctx, func(msg *sarama.ConsumerMessage) {}, wg)
	if err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func (a *App) Close() error {
	a.consumer.Close()
	return nil
}
