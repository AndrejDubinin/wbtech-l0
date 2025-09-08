package get

import (
	"context"
	"errors"
	"testing"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetOrder(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	repoMock := NewRepositoryMock(mc)
	cacheMock := NewCacheMock(mc)

	getOrderUsecase := New(repoMock, cacheMock)

	errDBDown := errors.New("db is down")
	order := domain.Order{
		OrderUID: "1",
	}

	tests := []struct {
		name     string
		orderUID string
		cache    *domain.Order
		dbErr    error
		want     *domain.Order
		wantErr  error
	}{
		{
			name:     "cache_order",
			orderUID: "1",
			cache:    &order,
			want:     &order,
			wantErr:  nil,
		},
		{
			name:     "db_order",
			orderUID: "1",
			cache:    nil,
			want:     &order,
			wantErr:  nil,
		},
		{
			name:     "db_error",
			orderUID: "1",
			cache:    nil,
			dbErr:    errDBDown,
			want:     nil,
			wantErr:  errDBDown,
		},
		{
			name:     "db_no_order",
			orderUID: "1",
			cache:    nil,
			dbErr:    nil,
			want:     nil,
			wantErr:  domain.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			cacheMock.GetMock.Expect(tt.orderUID).Return(tt.cache)
			if tt.cache == nil {
				repoMock.GetOrderMock.Expect(ctx, tt.orderUID).Return(tt.want, tt.dbErr)
			}

			got, err := getOrderUsecase.GetOrder(ctx, tt.orderUID)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
