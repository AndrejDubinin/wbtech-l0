package http

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/AndrejDubinin/wbtech-l0/internal/app/definitions"
	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type (
	getOrderUsecase interface {
		GetOrder(ctx context.Context, orderUID string) (*domain.Order, error)
	}

	GetOrderHandler struct {
		name            string
		getOrderUsecase getOrderUsecase
	}
)

var (
	re                     = regexp.MustCompile("^[a-zA-Z0-9_-]{8,64}$")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrInternalServerError = errors.New("internal server error")
)

func NewGetOrderHandler(usecase getOrderUsecase, name string) *GetOrderHandler {
	return &GetOrderHandler{
		name:            name,
		getOrderUsecase: usecase,
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
			// TODO: add place like "const op = app.http.GetOrderHandler" to log
			log.Printf("order not found, orderUID: %s\n", orderUID)
			GetErrorResponse(w, http.StatusNotFound, err, "")
			return
		}
		log.Printf("getOrderUsecase.GetOrder: %v\n", err)
		GetErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError, "")
		return
	}

	response, err := json.Marshal(order)
	if err != nil {
		log.Printf("get order json.Marshal: %v\n", err)
		GetErrorResponse(w, http.StatusInternalServerError, ErrInternalServerError, "")
		return
	}
	GetSuccessResponseWithBody(w, response)
}
