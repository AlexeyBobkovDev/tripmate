package generate_route_postgres_repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/repository/postgres/pool/mocks"
	generate_route_service "github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/features/generate-route/service"
)

func TestGenerateRouteRepository_FindCandidates(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRows := mocks.NewMockRows(ctrl)
	mockPool := mocks.NewMockPool(ctrl)
	routesRepository := NewGenerateRouteRepository(mockPool)

	mockPool.EXPECT().
		Query(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockRows, nil)
	mockRows.EXPECT().Next().Return(false)
	mockRows.EXPECT().Err().Return(nil)
	mockRows.EXPECT().Close()

	candidates, err := routesRepository.FindCandidates(context.Background(), generate_route_service.FindCandidatesCriteria{})

	require.NoError(t, err)
	assert.Empty(t, candidates)
}
