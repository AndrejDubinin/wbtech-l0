package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type (
	addOrderUsecase interface {
		AddOrder(ctx context.Context, order domain.Order) error
	}

	Handler struct {
		validate        *validator.Validate
		ServeMsgFn      func(context.Context, *sarama.ConsumerMessage)
		addOrderUsecase addOrderUsecase
	}
)

func NewHandler(usecase addOrderUsecase) *Handler {
	handler := &Handler{
		validate:        validator.New(validator.WithRequiredStructEnabled()),
		addOrderUsecase: usecase,
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
		log.Printf("consumer.handler decode message: %v", err)
		return
	}

	if err := h.validate.Struct(order); err != nil {
		log.Printf("consumer.handler validation: %v", err)
		return
	}

	err := h.addOrderUsecase.AddOrder(ctx, order)
	if err != nil {
		log.Printf("consumer.handler usecase: %v", err)
	}
}
