package routes_postgres_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/mocks"
)

func TestRoutesRepository_GetRoute_ReturnsNotFoundOnNoRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	tripID, routeID := 1, 42

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(core_postgres_pool_errors.ErrNoRows)
	mockPool.EXPECT().QueryRow(gomock.Any(), getRouteQuery, tripID, routeID).Return(mockRow)

	route, err := routesRepository.GetRoute(context.Background(), tripID, routeID)

	if route != nil {
		t.Fatal("route must be nil on no rows")
	}
	if !errors.Is(err, core_errors.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestRoutesRepository_GetRoute_WrapsOtherScanErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockPool := mocks.NewMockPool(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	routesRepository := NewRoutesRepository(mockPool)

	mockPool.EXPECT().OpTimeout().Return(5 * time.Second)
	mockRow.EXPECT().
		Scan(
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		Return(errors.New("connection reset"))
	mockPool.EXPECT().QueryRow(gomock.Any(), getRouteQuery, 1, 1).Return(mockRow)

	_, err := routesRepository.GetRoute(context.Background(), 1, 1)

	if err == nil || errors.Is(err, core_errors.ErrNotFound) {
		t.Fatalf("err = %v, want a wrapped non-NotFound error", err)
	}
}
