package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka"
)

type (
	Config struct {
		Topic string
	}
	Handler interface {
		ServeMsg(context.Context, *sarama.ConsumerMessage)
	}
	logger interface {
		Info(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
	}

	Consumer struct {
		config   Config
		consumer sarama.Consumer
		logger   logger
	}
)

func NewConsumer(kafkaConfig kafka.Config, conf Config, logger logger, opts ...Option) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	for _, opt := range opts {
		err := opt.Apply(config)
		if err != nil {
			return nil, err
		}
	}

	consumer, err := sarama.NewConsumer(kafkaConfig.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		config:   conf,
		logger:   logger,
	}, err
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) ConsumeTopic(ctx context.Context, handler Handler, wg *sync.WaitGroup) error {
	partitionList, err := c.consumer.Partitions(c.config.Topic)
	if err != nil {
		return err
	}

	initialOffset := sarama.OffsetNewest

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(c.config.Topic, partition, initialOffset)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func(pc sarama.PartitionConsumer, partition int32) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					c.logger.Info("consumer terminated",
						zap.String("topic", c.config.Topic),
						zap.Int32("partition", partition))
					return
				case msg, ok := <-pc.Messages():
					if !ok {
						c.logger.Info("consumer mag channel closed")
						return
					}
					handler.ServeMsg(ctx, msg)
				}
			}
		}(pc, partition)
	}

	return nil
}
