package app

import (
	"fmt"

	"github.com/AndrejDubinin/wbtech-l0/internal/app/definitions"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka"
	"github.com/AndrejDubinin/wbtech-l0/internal/infra/kafka/consumer"
)

type (
	Options struct {
		KafkaBrokerAddr, KafkaTopicName, DBConnStr, Addr string
		CacheCapacity                                    int64
	}
	path struct {
		index, orderItemGet string
	}

	config struct {
		kafka         kafka.Config
		consumer      consumer.Config
		dbConnStr     string
		cacheCapacity int64
		addr          string
		path          path
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
		addr:          opts.Addr,
		path: path{
			index:        "/",
			orderItemGet: fmt.Sprintf("/order/{%s}", definitions.ParamOrderUID),
		},
	}
}
