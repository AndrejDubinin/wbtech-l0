package add

import (
	"context"
	"fmt"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type (
	repository interface {
		AddOrder(ctx context.Context, order domain.Order) error
	}
	cache interface {
		Put(order *domain.Order)
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

func (u *Usecase) AddOrder(ctx context.Context, order domain.Order) error {
	err := u.repo.AddOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("repo.AddOrder: %v", err)
	}

	u.cache.Put(&order)

	return nil
}
