package routes_postgres_repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

const getRoutesQuery = `
SELECT id,
	trip_id,
	trip_version,
	version,
	status,
	created_at,
	updated_at
FROM app.routes
WHERE trip_id=$1
ORDER BY id;
`

func (r *RoutesRepository) GetRoutes(
	ctx context.Context,
	tripID int,
) ([]*domain.Route, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var routeResponses []*domain.Route

	routeRows, err := r.pool.Query(
		ctx,
		getRoutesQuery,
		tripID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route rows: %w", err)
	}
	defer routeRows.Close()

	for routeRows.Next() {
		var routeResponse domain.Route
		if err := routeRows.Scan(
			&routeResponse.ID,
			&routeResponse.TripID,
			&routeResponse.TripVersion,
			&routeResponse.Version,
			&routeResponse.Status,
			&routeResponse.CreatedAt,
			&routeResponse.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to fetch route: %w", err)
		}
		routeResponses = append(routeResponses, &routeResponse)
	}
	if err := routeRows.Err(); err != nil {
		return nil, fmt.Errorf("error while fetching routes rows: %w", err)
	}

	for routeRowsInd := range routeResponses {
		route, err := getRoute(ctx, r.pool, tripID, routeResponses[routeRowsInd].ID)
		if err != nil {
			return nil, err
		}
		routeResponses[routeRowsInd] = route
	}

	return routeResponses, nil
}
