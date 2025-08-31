package consumer

import (
	"context"
	"slices"
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

func NewConsumer(ctx context.Context, kafkaConfig kafka.Config, conf Config, logger logger, opts ...Option) (*Consumer, error) {
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

	ctxTopic, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	if err := waitForTopic(ctxTopic, kafkaConfig.Brokers, conf.Topic, logger); err != nil {
		return nil, err
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

func waitForTopic(ctx context.Context, brokers []string, topic string, logger logger) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			config := sarama.NewConfig()
			client, err := sarama.NewClient(brokers, config)
			if err != nil {
				logger.Error("kafka not ready", zap.Error(err))
				continue
			}

			topics, err := client.Topics()
			if err := client.Close(); err != nil {
				logger.Error("sarama client close", zap.Error(err))
			}
			if err != nil {
				logger.Error("list topics", zap.Error(err))
				continue
			}

			if slices.Contains(topics, topic) {
				logger.Info("topic exists", zap.String("topic", topic))
				return nil
			}

			logger.Info("topic not found, retrying...", zap.String("topic", topic))
		}
	}
}
