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

	Usecase struct {
		repo repository
	}
)

func New(repo repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) AddOrder(ctx context.Context, order domain.Order) error {
	err := u.repo.AddOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("repo.AddOrder: %v", err)
	}
	return nil
}
