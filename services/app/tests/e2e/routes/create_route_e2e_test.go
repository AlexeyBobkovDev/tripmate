//go:build e2e

package routes_e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
)

func TestRoutesTransportHTTP_CreateRoute_E2E_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Create route E2E test")

	createdRoute := createRoute(t, ctx, user.ID, trip.ID)

	generatedRoute := expectedGeneratedRoute(trip)
	assertRouteResponseMatchesGenerated(t, generatedRoute, createdRoute, trip)

	retrievedRoute, statusCode := getRoute(t, ctx, trip.ID, createdRoute.ID)
	require.Equal(t, http.StatusOK, statusCode)
	require.NotNil(t, retrievedRoute)

	assert.Equal(t, createdRoute.ID, retrievedRoute.ID)
	assertRouteResponseMatchesGenerated(t, generatedRoute, retrievedRoute, trip)
}

func TestRoutesTransportHTTP_CreateRoute_E2E_InvalidOwnerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Create route invalid owner E2E test")

	invalidOwnerIDs := []string{"abc", "0", "-1"}

	for _, ownerID := range invalidOwnerIDs {
		t.Run(ownerID, func(t *testing.T) {
			response := doRequest(
				t,
				ctx,
				http.MethodPost,
				fmt.Sprintf("%s/routes/%s/%d", baseURL, ownerID, trip.ID),
				nil,
			)
			defer response.Body.Close()

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestRoutesTransportHTTP_CreateRoute_E2E_InvalidTripID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	invalidTripIDs := []string{"abc", "0", "-1"}

	for _, tripID := range invalidTripIDs {
		t.Run(tripID, func(t *testing.T) {
			response := doRequest(
				t,
				ctx,
				http.MethodPost,
				fmt.Sprintf("%s/routes/%d/%s", baseURL, user.ID, tripID),
				nil,
			)
			defer response.Body.Close()

			assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		})
	}
}

func TestRoutesTransportHTTP_CreateRoute_E2E_TripNotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)

	response := doRequest(
		t,
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/routes/%d/%d", baseURL, user.ID, 2_000_000_000),
		nil,
	)
	defer response.Body.Close()

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestRoutesTransportHTTP_CreateRoute_E2E_WrongOwner(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	owner := createTestUser(t, ctx)
	anotherUser := createTestUser(t, ctx)
	trip := createTrip(t, ctx, owner.ID, "Create route wrong owner E2E test")

	response := doRequest(
		t,
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/routes/%d/%d", baseURL, anotherUser.ID, trip.ID),
		nil,
	)
	defer response.Body.Close()

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestRoutesTransportHTTP_CreateRoute_E2E_GeneratorFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Create route generator failure E2E test")

	setGenerateRouteError(t, errors.New("route generator is unavailable"))

	response := doRequest(
		t,
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/routes/%d/%d", baseURL, user.ID, trip.ID),
		nil,
	)
	defer response.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestRoutesTransportHTTP_CreateRoute_E2E_GeneratorReturnsInvalidRoute(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	user := createTestUser(t, ctx)
	trip := createTrip(t, ctx, user.ID, "Create route invalid generated route E2E test")

	restore := testRouteGeneratorServer.setGenerateRouteFunc(func(
		ctx context.Context,
		request *route_generator_api_v1.GenerateRouteRequest,
	) (*route_generator_api_v1.GenerateRouteResponse, error) {
		return &route_generator_api_v1.GenerateRouteResponse{}, nil
	})
	t.Cleanup(restore)

	response := doRequest(
		t,
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/routes/%d/%d", baseURL, user.ID, trip.ID),
		nil,
	)
	defer response.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	routes, statusCode := getRoutes(t, ctx, trip.ID)
	require.Equal(t, http.StatusOK, statusCode)
	assert.Empty(t, routes)
}
