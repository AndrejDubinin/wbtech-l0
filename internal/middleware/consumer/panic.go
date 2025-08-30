package consumer

import (
	"context"
	"log"

	"github.com/IBM/sarama"

	"github.com/AndrejDubinin/wbtech-l0/internal/app/consumer"
)

func Panic(next *consumer.Handler) *consumer.Handler {
	return &consumer.Handler{
		ServeMsgFn: func(ctx context.Context, msg *sarama.ConsumerMessage) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic recovered in consumer: %v (topic=%s, partition=%d, offset=%d)", r, msg.Topic, msg.Partition, msg.Offset)
				}
			}()
			next.ServeMsgFn(ctx, msg)
		},
	}
}
