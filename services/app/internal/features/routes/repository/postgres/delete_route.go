package routes_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
)

const deleteRouteQuery = `
DELETE FROM app.routes
WHERE trip_id=$1 AND id=$2;
`

func (r *RoutesRepository) DeleteRoute(
	ctx context.Context,
	tripID int,
	routeID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	cmdTag, err := tx.Exec(
		ctx,
		deleteRouteQuery,
		tripID,
		routeID,
	)
	if err != nil {
		return fmt.Errorf("delete route row: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("the route %d does not exist: %w", routeID, core_errors.ErrNotFound)
	}
	return tx.Commit(ctx)
}
