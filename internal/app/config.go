package app

import (
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka/consumer"
)

type (
	Options struct {
		KafkaBrokerAddr, KafkaTopicName, DBConnStr string
		CacheCapacity                              int64
	}

	config struct {
		kafka         kafka.Config
		consumer      consumer.Config
		dbConnStr     string
		cacheCapacity int64
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
		dbConnStr:     opts.DBConnStr,
		cacheCapacity: opts.CacheCapacity,
	}
}
