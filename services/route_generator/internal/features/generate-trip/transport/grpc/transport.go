package transport_grpc_generate_route

import (
	"context"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	"github.com/AlexeyBobkovDev/tripmate/services/route_generator/internal/core/domain"
)

type GenerateRouteService interface {
	GenerateRoute(
		ctx context.Context,
		trip *domain.Trip,
	) (*domain.Route, error)
}
type GenerateRouteGRPCTransport struct {
	route_generator_api_v1.UnimplementedRouteGeneratorServiceServer
	generateRouteService GenerateRouteService
}

func NewGenerateRouteServer(
	generateRouteService GenerateRouteService,
) *GenerateRouteGRPCTransport {
	return &GenerateRouteGRPCTransport{
		generateRouteService: generateRouteService,
	}
}
