package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka"
	"github.com/IBM/sarama"
)

type Config struct {
	Topic string
}

type Consumer struct {
	config   Config
	consumer sarama.Consumer
}

func NewConsumer(kafkaConfig kafka.Config, conf Config, opts ...Option) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	for _, opt := range opts {
		_ = opt.Apply(config)
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

func (c *Consumer) ConsumeTopic(ctx context.Context, handler func(*sarama.ConsumerMessage), wg *sync.WaitGroup) error {
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
					fmt.Printf("consumer for topic=%s, partition=%d terminated\n", c.config.Topic, partition)
					return
				case msg, ok := <-pc.Messages():
					if !ok {
						fmt.Println("consumer mag channel closed")
						return
					}
					handler(msg)
				}
			}
		}(pc, partition)
	}

	return nil
}
