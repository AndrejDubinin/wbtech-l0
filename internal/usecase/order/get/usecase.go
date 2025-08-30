package get

import (
	"context"
	"fmt"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type (
	repository interface {
		GetOrder(ctx context.Context, orderUID string) (*domain.Order, error)
	}
	cache interface {
		Get(orderUID string) *domain.Order
	}

	Usecase struct {
		repo  repository
		cache cache
	}
)

func New(repo repository, cache cache) *Usecase {
	return &Usecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *Usecase) GetOrder(ctx context.Context, orderUID string) (*domain.Order, error) {
	order := u.cache.Get(orderUID)
	if order != nil {
		return order, nil
	}

	order, err := u.repo.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetOrder: %v", err)
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	return order, nil
}
