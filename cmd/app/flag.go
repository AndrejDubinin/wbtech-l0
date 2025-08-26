package main

import (
	"flag"
	"fmt"

	"github.com/AndrejDubinin/wbtech-l0/internal/app"
)

const (
	defaultKafkaBrokerAddr = "localhost:9092"
	defaultKafkaTopicName  = "wbtech-l0-topic"
)

var opts = app.Options{}

func initOpts() {
	flag.StringVar(&opts.KafkaBrokerAddr, "broker_addr", defaultKafkaBrokerAddr, fmt.Sprintf("kafka broker host and port, default: %q", defaultKafkaBrokerAddr))
	flag.StringVar(&opts.KafkaTopicName, "topic_name", defaultKafkaTopicName, fmt.Sprintf("kafka topic's name, default: %q", defaultKafkaTopicName))
	flag.Parse()
}
