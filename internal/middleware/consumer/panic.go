package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/AndrejDubinin/wbtech-l0/internal/app/consumer"
)

type (
	logger interface {
		Info(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
	}
)

func Panic(next *consumer.Handler, logger logger) *consumer.Handler {
	return &consumer.Handler{
		ServeMsgFn: func(ctx context.Context, msg *sarama.ConsumerMessage) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("panic recovered in consumer",
						zap.Any("error", r),
						zap.String("topic", msg.Topic),
						zap.Int32("partition", msg.Partition),
						zap.Int64("offset", msg.Offset))
				}
			}()
			next.ServeMsgFn(ctx, msg)
		},
	}
}
