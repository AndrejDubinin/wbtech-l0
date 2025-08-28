package app

import (
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka/consumer"
)

type (
	Options struct {
		KafkaBrokerAddr, KafkaTopicName, DbConnStr string
	}

	config struct {
		kafka     kafka.Config
		consumer  consumer.Config
		dbConnStr string
	}
)

func NewConfig(opts Options) config {
	return config{
		kafka: kafka.Config{
			Brokers: []string{
				opts.KafkaBrokerAddr,
			},
		},
		consumer: consumer.Config{
			Topic: opts.KafkaTopicName,
		},
		dbConnStr: opts.DbConnStr,
	}
}
