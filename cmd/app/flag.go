package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/AndrejDubinin/wbtech-l0/internal/app"
)

const (
	defaultKafkaBrokerAddr = "localhost:9092"
	defaultKafkaTopicName  = "wbtech-l0-topic"

	dbConnStr     = "DB_CONN"
	cacheCapacity = "CACHE_CAPACITY"
)

var opts = app.Options{}

func initOpts() {
	flag.StringVar(&opts.KafkaBrokerAddr, "broker_addr", defaultKafkaBrokerAddr, fmt.Sprintf("kafka broker host and port, default: %q", defaultKafkaBrokerAddr))
	flag.StringVar(&opts.KafkaTopicName, "topic_name", defaultKafkaTopicName, fmt.Sprintf("kafka topic's name, default: %q", defaultKafkaTopicName))
	flag.Parse()

	opts.DBConnStr = os.Getenv(dbConnStr)
	cacheCapacity, err := strconv.ParseInt(os.Getenv(cacheCapacity), 10, 64)
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}
	opts.CacheCapacity = cacheCapacity
}
