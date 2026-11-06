package transport_grpc_generate_route

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
)

func (r *GenerateRouteGRPCTransport) GenerateRoute(
	ctx context.Context,
	in *route_generator_api_v1.GenerateRouteRequest,
) (*route_generator_api_v1.GenerateRouteResponse, error) {
	generatedRoute, err := r.generateRouteService.GenerateRoute(ctx, convertTripDTOToDomain(in))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertRouteDomainToDTO(generatedRoute), nil
}
