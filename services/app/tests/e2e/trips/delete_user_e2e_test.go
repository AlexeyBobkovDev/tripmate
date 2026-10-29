//go:build e2e

package trips_e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersTransportHTTP_DeleteUser_E2E_SoftDeletesUserTrips(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Trip To Soft Delete On User Delete")

	deleteUserRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/users/%d", baseURL, user.ID),
		nil,
	)
	require.NoError(t, err)

	response, err := newHTTPClient().Do(deleteUserRequest)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusNoContent, response.StatusCode)

	var tripDeletedAt time.Time
	require.NoError(t, testDeps.pool.QueryRow(
		ctx,
		`SELECT deleted_at FROM app.trips WHERE owner_id=$1 AND id=$2`,
		user.ID,
		trip.ID,
	).Scan(&tripDeletedAt))
	assert.False(t, tripDeletedAt.IsZero())

	var userDeletedAt time.Time
	require.NoError(t, testDeps.pool.QueryRow(
		ctx,
		`SELECT deleted_at FROM app.users WHERE id=$1`,
		user.ID,
	).Scan(&userDeletedAt))
	assert.False(t, userDeletedAt.IsZero())
}
