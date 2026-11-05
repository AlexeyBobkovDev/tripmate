//go:build e2e

package routes_e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	route_generator_api_v1 "github.com/AlexeyBobkovDev/tripmate/contracts/gen/go/tripmate/route_generator/api/v1"
	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	routes_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/routes/transport/http"
	trips_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

var (
	baseURL          = fmt.Sprintf("http://localhost%s/api/v1", freePort)
	requestTimeout   = 5 * time.Second
	testUserSequence atomic.Int64
)

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: requestTimeout}
}

func registerCleanup(t *testing.T, cleanup func(context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := cleanup(ctx); err != nil {
			t.Errorf("E2E cleanup failed: %v", err)
		}
	})
}

func doRequest(
	t *testing.T,
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, method, path, body)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)

	return response
}

func createTestUser(t *testing.T, ctx context.Context) *users_transport_http.UserResponse {
	t.Helper()

	sequence := testUserSequence.Add(1)
	suffix := fmt.Sprintf("%d_%d", time.Now().UnixNano(), sequence)

	body := fmt.Sprintf(`{
		"name": %q,
		"surname": %q,
		"username": %q,
		"birth_date": %q,
		"description": %q,
		"email": %q,
		"phone_number": %q,
		"password": %q
	}`,
		testutils.TestName,
		testutils.TestSurname,
		"routes_user_"+suffix,
		testutils.BirthDateTest().Format(core_consts.Layout),
		testutils.TestDescription,
		"routes_"+suffix+"@example.com",
		fmt.Sprintf("+7999%07d", sequence%10_000_000),
		testutils.TestPassword,
	)

	response := doRequest(
		t,
		ctx,
		http.MethodPost,
		baseURL+"/users",
		strings.NewReader(body),
	)
	defer response.Body.Close()

	require.Equal(t, http.StatusCreated, response.StatusCode)

	var user users_transport_http.UserResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&user))
	require.NotZero(t, user.ID)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.users WHERE id = $1`,
			user.ID,
		)
		return err
	})

	return &user
}

type tripRequestCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type tripRequestBudget struct {
	AmountMinorUnits int64  `json:"amount_minor_units"`
	Currency         string `json:"currency"`
}

type tripRequestBody struct {
	OwnerID          int                    `json:"owner_id"`
	Version          int                    `json:"version,omitempty"`
	Name             string                 `json:"name"`
	PlanningArea     string                 `json:"planning_area"`
	StartCoordinates tripRequestCoordinates `json:"start_location"`
	StartDate        string                 `json:"start_date"`
	EndDate          string                 `json:"end_date"`
	Budget           tripRequestBudget      `json:"budget"`
	Interests        []string               `json:"interests"`
	TimeZone         string                 `json:"time_zone"`
}

func newTripRequestBody(ownerID int, name string) tripRequestBody {
	startDate := time.Now().Add(48 * time.Hour)
	endDate := startDate

	return tripRequestBody{
		OwnerID:      ownerID,
		Name:         name,
		PlanningArea: "Moscow",
		StartCoordinates: tripRequestCoordinates{
			Latitude:  55.7539,
			Longitude: 37.6208,
		},
		StartDate: startDate.Format(core_consts.Layout),
		EndDate:   endDate.Format(core_consts.Layout),
		Budget: tripRequestBudget{
			AmountMinorUnits: 1500000,
			Currency:         "RUB",
		},
		Interests: []string{"history", "nature"},
		TimeZone:  "Europe/Moscow",
	}
}

func marshalJSON(t *testing.T, value any) *bytes.Reader {
	t.Helper()

	body, err := json.Marshal(value)
	require.NoError(t, err)

	return bytes.NewReader(body)
}

func createTrip(
	t *testing.T,
	ctx context.Context,
	ownerID int,
	name string,
) *trips_transport_http.TripResponse {
	t.Helper()

	response := doRequest(
		t,
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/trips/%d", baseURL, ownerID),
		marshalJSON(t, newTripRequestBody(ownerID, name)),
	)
	defer response.Body.Close()

	require.Equal(t, http.StatusCreated, response.StatusCode)

	var trip trips_transport_http.TripResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&trip))
	require.NotZero(t, trip.ID)
	require.Equal(t, ownerID, trip.OwnerID)

	registerCleanup(t, func(ctx context.Context) error {
		_, err := testDeps.pool.Exec(
			ctx,
			`DELETE FROM app.trips WHERE id = $1`,
			trip.ID,
		)
		return err
	})

	return &trip
}

func createRoute(
	t *testing.T,
	ctx context.Context,
	ownerID int,
	tripID int,
) *routes_transport_http.RouteResponse {
	t.Helper()

	response := doRequest(
		t,
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/routes/%d/%d", baseURL, ownerID, tripID),
		nil,
	)
	defer response.Body.Close()

	require.Equal(t, http.StatusCreated, response.StatusCode)

	var route routes_transport_http.RouteResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&route))
	require.NotZero(t, route.ID)
	require.Equal(t, tripID, route.TripID)

	return &route
}

func getRoute(
	t *testing.T,
	ctx context.Context,
	tripID int,
	routeID int,
) (*routes_transport_http.RouteResponse, int) {
	t.Helper()

	response := doRequest(
		t,
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/routes/%d/%d", baseURL, tripID, routeID),
		nil,
	)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode
	}

	var route routes_transport_http.RouteResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&route))

	return &route, response.StatusCode
}

func getRoutes(
	t *testing.T,
	ctx context.Context,
	tripID int,
) ([]routes_transport_http.RouteResponse, int) {
	t.Helper()

	response := doRequest(
		t,
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/routes/%d", baseURL, tripID),
		nil,
	)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode
	}

	var routes []routes_transport_http.RouteResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&routes))

	return routes, response.StatusCode
}

func deleteRoute(
	t *testing.T,
	ctx context.Context,
	tripID int,
	routeID int,
) int {
	t.Helper()

	response := doRequest(
		t,
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/routes/%d/%d", baseURL, tripID, routeID),
		nil,
	)
	defer response.Body.Close()

	return response.StatusCode
}

func decodeErrorResponse(
	t *testing.T,
	response *http.Response,
) core_http_response.ErrorResponse {
	t.Helper()

	var result core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&result))

	return result
}

func setGenerateRouteError(t *testing.T, err error) {
	t.Helper()

	restore := testRouteGeneratorServer.setGenerateRouteFunc(func(
		ctx context.Context,
		request *route_generator_api_v1.GenerateRouteRequest,
	) (*route_generator_api_v1.GenerateRouteResponse, error) {
		return nil, err
	})
	t.Cleanup(restore)
}

func expectedGeneratedRoute(trip *trips_transport_http.TripResponse) *route_generator_api_v1.GenerateRouteResponse {
	return testutils.RouteGeneratorResponseTest(trip.StartDate)
}

func assertRouteResponseMatchesGenerated(
	t *testing.T,
	generated *route_generator_api_v1.GenerateRouteResponse,
	actual *routes_transport_http.RouteResponse,
	trip *trips_transport_http.TripResponse,
) {
	t.Helper()

	require.NotNil(t, actual)

	assert.NotZero(t, actual.ID)
	assert.Equal(t, trip.ID, actual.TripID)
	assert.Equal(t, trip.Version, actual.TripVersion)
	assert.Equal(t, 1, actual.Version)
	assert.Equal(t, domain.RouteStatus(domain.RouteStatusGenerated), actual.Status)
	assert.False(t, actual.CreatedAt.IsZero())
	assert.False(t, actual.UpdatedAt.IsZero())

	require.Len(t, actual.Days, len(generated.Days))

	for _, generatedDay := range generated.Days {
		actualDay, found := routeDayByNumber(actual.Days, int(generatedDay.DayNumber))
		require.Truef(t, found, "route day %d is not found", generatedDay.DayNumber)

		assert.NotZero(t, actualDay.ID)
		assertSameDate(t, generatedDay.Date.AsTime(), actualDay.Date)
		assert.Equal(t, int(generatedDay.DayNumber), actualDay.DayNumber)
		require.Len(t, actualDay.Stops, len(generatedDay.Stops))

		stopIDsMapping := make(map[int64]int, len(generatedDay.Stops))

		for _, generatedStop := range generatedDay.Stops {
			actualStop, stopFound := routeStopByPosition(actualDay, int(generatedStop.Position))
			require.Truef(
				t,
				stopFound,
				"route stop with position %d is not found in day %d",
				generatedStop.Position,
				generatedDay.DayNumber,
			)

			assert.NotZero(t, actualStop.ID)
			assert.Equal(t, int(generatedStop.PlaceId), actualStop.PlaceID)
			assert.Equal(t, int(generatedStop.Position), actualStop.Position)
			assert.Equal(t, generatedStop.Name, actualStop.Name)
			assert.Equal(t, generatedStop.Description, actualStop.Description)
			assert.Equal(t, generatedStop.Coordinates.Latitude, actualStop.Coordinates.Latitude)
			assert.Equal(t, generatedStop.Coordinates.Longitude, actualStop.Coordinates.Longitude)
			assertSameInstant(t, generatedStop.ArrivalTime.AsTime(), actualStop.ArrivalTime)
			assertSameInstant(t, generatedStop.DepartureTime.AsTime(), actualStop.DepartureTime)
			assert.Equal(t, generatedStop.VisitDuration.AsDuration(), actualStop.VisitDuration)

			stopIDsMapping[generatedStop.Id] = actualStop.ID
		}

		require.Len(t, actualDay.Legs, len(generatedDay.Legs))

		for _, generatedLeg := range generatedDay.Legs {
			expectedFromStopID := stopIDsMapping[int64(generatedLeg.FromStopId)]
			expectedToStopID := stopIDsMapping[int64(generatedLeg.ToStopId)]

			actualLeg, legFound := routeLegByStopIDs(actualDay, expectedFromStopID, expectedToStopID)
			require.Truef(
				t,
				legFound,
				"route leg %d -> %d is not found in day %d",
				expectedFromStopID,
				expectedToStopID,
				generatedDay.DayNumber,
			)

			assert.NotZero(t, actualLeg.ID)
			assert.Equal(t, expectedFromStopID, actualLeg.FromStopID)
			assert.Equal(t, expectedToStopID, actualLeg.ToStopID)
			assert.Equal(t, domain.TravelModeWalking, actualLeg.TravelMode)
			assert.Equal(t, int(generatedLeg.Distance), actualLeg.Distance)
			assert.Equal(t, generatedLeg.Duration.AsDuration(), actualLeg.Duration)
		}
	}
}

func routeDayByNumber(days []domain.RouteDay, dayNumber int) (domain.RouteDay, bool) {
	for _, day := range days {
		if day.DayNumber == dayNumber {
			return day, true
		}
	}

	return domain.RouteDay{}, false
}

func routeStopByPosition(day domain.RouteDay, position int) (domain.RouteStop, bool) {
	for _, stop := range day.Stops {
		if stop.Position == position {
			return stop, true
		}
	}

	return domain.RouteStop{}, false
}

func routeLegByStopIDs(
	day domain.RouteDay,
	fromStopID int,
	toStopID int,
) (domain.RouteLeg, bool) {
	for _, leg := range day.Legs {
		if leg.FromStopID == fromStopID && leg.ToStopID == toStopID {
			return leg, true
		}
	}

	return domain.RouteLeg{}, false
}

func assertSameDate(t *testing.T, expected time.Time, actual time.Time) {
	t.Helper()

	assert.Equal(t, expected.Format("2006-01-02"), actual.Format("2006-01-02"))
}

func assertSameInstant(t *testing.T, expected time.Time, actual time.Time) {
	t.Helper()

	assert.Truef(
		t,
		expected.Equal(actual),
		"expected time %s, got %s",
		expected.Format(time.RFC3339Nano),
		actual.Format(time.RFC3339Nano),
	)
}
