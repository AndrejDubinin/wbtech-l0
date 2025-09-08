package add

import (
	"context"
	"errors"
	"testing"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestAddOrder(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	repoMock := NewRepositoryMock(mc)
	cacheMock := NewCacheMock(mc)

	addOrderUsecase := New(repoMock, cacheMock)
	errDBDown := errors.New("db is down")
	order := domain.Order{
		OrderUID: "1",
	}

	tests := []struct {
		name    string
		order   domain.Order
		wantErr error
	}{
		{
			name:    "db_error",
			order:   order,
			wantErr: errDBDown,
		},
		{
			name:    "db_order",
			order:   order,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repoMock.AddOrderMock.Expect(ctx, tt.order).Return(tt.wantErr)
			if tt.wantErr == nil {
				cacheMock.PutMock.Expect(&tt.order).Return()
			}

			err := addOrderUsecase.AddOrder(ctx, tt.order)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
