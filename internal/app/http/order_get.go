package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/wbtech-l0/internal/app/definitions"
	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type (
	getOrderUsecase interface {
		GetOrder(ctx context.Context, orderUID string) (*domain.Order, error)
	}
	logger interface {
		Info(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
	}

	GetOrderHandler struct {
		name            string
		getOrderUsecase getOrderUsecase
		logger          logger
	}
)

var (
	re                     = regexp.MustCompile("^[a-zA-Z0-9_-]{8,64}$")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrInternalServerError = errors.New("internal server error")
)

func NewGetOrderHandler(usecase getOrderUsecase, name string, logger logger) *GetOrderHandler {
	return &GetOrderHandler{
		name:            name,
		getOrderUsecase: usecase,
		logger:          logger,
	}
}

func (h *GetOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orderUID := r.PathValue(definitions.ParamOrderUID)
	if !re.MatchString(orderUID) {
		GetErrorResponse(w, http.StatusBadRequest, ErrInvalidParameter, "")
		return
	}

	ctx := r.Context()
	order, err := h.getOrderUsecase.GetOrder(ctx, orderUID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			h.logger.Error("order not found", zap.String("orderUID", orderUID))
			GetErrorResponse(w, http.StatusNotFound, err, "")
			return
		}
		h.logger.Error("getOrderUsecase.GetOrder", zap.Error(err))
		GetErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError, "")
		return
	}

	response, err := json.Marshal(order)
	if err != nil {
		h.logger.Error("json.Marshal", zap.Error(err))
		GetErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError, "")
		return
	}
	GetSuccessResponseWithBody(w, response)
}
