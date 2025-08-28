package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	appConsumer "github.com/AndrejDubinin/wbtech-l0/internal/app/consumer"
	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
	memoryorder "github.com/AndrejDubinin/wbtech-l0/internal/infra/cache/memory_order"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka/consumer"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/repository/order"
	consumerMw "github.com/AndrejDubinin/wbtech-l0/internal/middleware/consumer"
	"github.com/AndrejDubinin/wbtech-l0/internal/usecase/order/add"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	cons interface {
		ConsumeTopic(ctx context.Context, handler consumer.Handler, wg *sync.WaitGroup) error
		Close() error
	}
	orderStorage interface {
		AddOrder(ctx context.Context, order domain.Order) error
	}
	orderCache interface {
		Get(orderUID string) *domain.Order
		Add(order *domain.Order)
	}

	App struct {
		config   config
		consumer cons
		db       *pgxpool.Pool
		storage  orderStorage
		cache    orderCache
	}
)

func NewApp(config config) (*App, error) {
	cons, err := consumer.NewConsumer(config.kafka, config.consumer,
		consumer.WithReturnErrorsEnabled(true),
	)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.dbConnStr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to db: %w", err))
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &App{
		config:   config,
		consumer: cons,
		db:       db,
		storage:  order.NewRepository(db),
		cache:    memoryorder.New(),
	}, nil
}

func (a *App) Run() error {
	var (
		wg  = &sync.WaitGroup{}
		ctx = context.Background()
	)

	err := a.runConsumer(ctx, wg)
	if err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func (a *App) Close() error {
	a.consumer.Close()
	a.db.Close()
	return nil
}

func (a *App) runConsumer(ctx context.Context, wg *sync.WaitGroup) error {
	consumerHandler := appConsumer.NewHandler(add.New(a.storage, a.cache))
	consumerHandler = consumerMw.Panic(consumerHandler)

	log.Printf("consumer reads topic: %s\n", a.config.consumer.Topic)
	err := a.consumer.ConsumeTopic(ctx, consumerHandler, wg)
	if err != nil {
		return err
	}

	return nil
}
