package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka"
)

type (
	Config struct {
		Topic string
	}

	Handler interface {
		ServeMsg(context.Context, *sarama.ConsumerMessage)
	}

	Consumer struct {
		config   Config
		consumer sarama.Consumer
	}
)

func NewConsumer(kafkaConfig kafka.Config, conf Config, opts ...Option) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

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

	initialOffset := sarama.OffsetOldest

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
					log.Printf("consumer for topic=%s, partition=%d terminated\n", c.config.Topic, partition)
					return
				case msg, ok := <-pc.Messages():
					if !ok {
						log.Println("consumer mag channel closed")
						return
					}
					handler.ServeMsg(ctx, msg)
				}
			}
		}(pc, partition)
	}

	return nil
}
