package preload

import (
	"context"
	"fmt"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type (
	repository interface {
		GetOrders(ctx context.Context, amount int64) ([]*domain.Order, error)
	}
	cache interface {
		Put(order *domain.Order)
	}

	CachePreloader struct {
		capacity int64
		repo     repository
		cache    cache
	}
)

func New(capacity int64, repo repository, cache cache) *CachePreloader {
	return &CachePreloader{
		capacity: capacity,
		repo:     repo,
		cache:    cache,
	}
}

func (c *CachePreloader) Preload(ctx context.Context) error {
	orders, err := c.repo.GetOrders(ctx, c.capacity)
	if err != nil {
		return fmt.Errorf("repo.GetOrders: %w", err)
	}

	for _, order := range orders {
		c.cache.Put(order)
	}

	return nil
}
