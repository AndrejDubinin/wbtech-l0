package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	appConsumer "github.com/AndrejDubinin/wbtech-l0/internal/app/consumer"
	appHttp "github.com/AndrejDubinin/wbtech-l0/internal/app/http"
	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
	memoryorder "github.com/AndrejDubinin/wbtech-l0/internal/infra/cache/memory_order"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka/consumer"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/repository/order"
	consumerMw "github.com/AndrejDubinin/wbtech-l0/internal/middleware/consumer"
	httpMw "github.com/AndrejDubinin/wbtech-l0/internal/middleware/http"
	"github.com/AndrejDubinin/wbtech-l0/internal/usecase/cache/preload"
	"github.com/AndrejDubinin/wbtech-l0/internal/usecase/order/add"
	"github.com/AndrejDubinin/wbtech-l0/internal/usecase/order/get"
)

type (
	cons interface {
		ConsumeTopic(ctx context.Context, handler consumer.Handler, wg *sync.WaitGroup) error
		Close() error
	}
	orderStorage interface {
		AddOrder(ctx context.Context, order domain.Order) error
		GetOrders(ctx context.Context, amount int64) ([]*domain.Order, error)
		GetOrder(ctx context.Context, orderUID string) (*domain.Order, error)
	}
	orderCache interface {
		Get(orderUID string) *domain.Order
		Put(order *domain.Order)
	}
	server interface {
		ListenAndServe() error
		Close() error
		Shutdown(ctx context.Context) error
	}
	mux interface {
		Handle(pattern string, handler http.Handler)
	}
	logger interface {
		Info(msg string, fields ...zap.Field)
		Warn(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
		Fatal(msg string, fields ...zap.Field)
		Sync() error
	}

	App struct {
		config   config
		consumer cons
		db       *pgxpool.Pool
		storage  orderStorage
		cache    orderCache
		server   server
		mux      mux
		logger   logger
	}
)

func NewApp(ctx context.Context, config config, logger *zap.Logger) (*App, error) {
	cons, err := consumer.NewConsumer(config.kafka, config.consumer, logger,
		consumer.WithReturnErrorsEnabled(true),
	)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.New(ctx, config.dbConnStr)
	if err != nil {
		logger.Fatal("connection to db", zap.Error(err))
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	handler := httpMw.AccessLogMiddleware(mux, logger)
	handler = httpMw.PanicMiddleware(handler, logger)

	return &App{
		config:   config,
		consumer: cons,
		db:       db,
		storage:  order.NewRepository(db),
		cache:    memoryorder.New(config.cacheCapacity),
		mux:      mux,
		server: &http.Server{
			Addr:         config.addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger: logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}

	defer func() {
		if err := a.Close(ctx); err != nil {
			a.logger.Error("app.Close", zap.Error(err))
		}
	}()

	cachPreloader := preload.New(a.config.cacheCapacity, a.storage, a.cache)
	a.logger.Info("cash preloding")
	if err := cachPreloader.Preload(ctx); err != nil {
		return err
	}

	if err := a.runConsumer(ctx, wg); err != nil {
		return err
	}

	go func() {
		a.logger.Info("Starting server", zap.String("address", a.config.addr))
		err := a.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal("starting server", zap.Error(err))
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) Close(ctx context.Context) error {
	var errs []error

	a.logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("server.Shutdown: %w", err))
	}

	a.logger.Info("closing consumer")
	if err := a.consumer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("consumer.Close: %w", err))
	}

	a.logger.Info("closing database pool")
	a.db.Close()

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}

	a.logger.Info("all resources closed successfully")
	return nil
}

func (a *App) runConsumer(ctx context.Context, wg *sync.WaitGroup) error {
	consumerHandler := appConsumer.NewHandler(add.New(a.storage, a.cache), a.logger)
	consumerHandler = consumerMw.Panic(consumerHandler, a.logger)

	a.logger.Info("consumer reads topic", zap.String("topic", a.config.consumer.Topic))
	err := a.consumer.ConsumeTopic(ctx, consumerHandler, wg)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ListenAndServe() error {
	a.mux.Handle(a.config.path.index, appHttp.NewIndexHandler())
	a.mux.Handle(a.config.path.orderItemGet, appHttp.NewGetOrderHandler(get.New(a.storage, a.cache),
		a.config.path.orderItemGet, a.logger))

	return a.server.ListenAndServe()
}
