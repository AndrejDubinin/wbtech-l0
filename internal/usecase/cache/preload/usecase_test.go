package preload

import (
	"context"
	"errors"
	"testing"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

const cacheSize = 100

func TestPreloadCache(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	repoMock := NewRepositoryMock(mc)
	cacheMock := NewCacheMock(mc)

	cachePreloadUsecase := New(cacheSize, repoMock, cacheMock)

	tests := []struct {
		name    string
		orders  []*domain.Order
		wantErr error
	}{
		{
			name:    "db_no_orders",
			orders:  []*domain.Order{},
			wantErr: nil,
		},
		{
			name:    "db_error",
			orders:  []*domain.Order{},
			wantErr: errors.New("db is down"),
		},
		{
			name: "db_orders",
			orders: []*domain.Order{
				{
					OrderUID: "1",
				},
				{
					OrderUID: "2",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repoMock.GetOrdersMock.Expect(ctx, cacheSize).Return(tt.orders, tt.wantErr)
			for _, order := range tt.orders {
				cacheMock.PutMock.When(order)
			}

			err := cachePreloadUsecase.Preload(ctx)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
