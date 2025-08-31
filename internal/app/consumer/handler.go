package consumer

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type (
	addOrderUsecase interface {
		AddOrder(ctx context.Context, order domain.Order) error
	}
	logger interface {
		Info(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
	}

	Handler struct {
		validate        *validator.Validate
		ServeMsgFn      func(context.Context, *sarama.ConsumerMessage)
		addOrderUsecase addOrderUsecase
		logger          logger
	}
)

func NewHandler(usecase addOrderUsecase, logger logger) *Handler {
	handler := &Handler{
		validate:        validator.New(validator.WithRequiredStructEnabled()),
		addOrderUsecase: usecase,
		logger:          logger,
	}
	handler.ServeMsgFn = handler.serveMsg

	return handler
}

func (h *Handler) ServeMsg(ctx context.Context, s *sarama.ConsumerMessage) {
	h.ServeMsgFn(ctx, s)
}

func (h *Handler) serveMsg(ctx context.Context, s *sarama.ConsumerMessage) {
	order := domain.Order{}
	if err := json.Unmarshal(s.Value, &order); err != nil {
		h.logger.Error("json.unmarshal", zap.Error(err))
		return
	}

	if err := h.validate.Struct(order); err != nil {
		h.logger.Error("validation", zap.Error(err))
		return
	}

	err := h.addOrderUsecase.AddOrder(ctx, order)
	if err != nil {
		h.logger.Error("usecase", zap.Error(err))
	}
}
