//go:build e2e

package trips_e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	trips_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/trips/transport/http"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

var (
	baseURL        = fmt.Sprintf("http://localhost%s/api/v1", freePort)
	requestTimeout = 5 * time.Second
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

func createTestUser(t *testing.T, ctx context.Context) *users_transport_http.UserResponse {
	t.Helper()

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
		testutils.TestUsername,
		testutils.BirthDateTest().Format(core_consts.Layout),
		testutils.TestDescription,
		testutils.TestEmail,
		testutils.TestPhoneNumber,
		testutils.TestPassword,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/users",
		strings.NewReader(body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
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

type tripRequestLocation struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type tripRequestBudget struct {
	AmountMinorUnits int64  `json:"amount_minor_units"`
	Currency         string `json:"currency"`
}

type tripRequestBody struct {
	OwnerID       int                 `json:"owner_id"`
	Version       int                 `json:"version,omitempty"`
	Name          string              `json:"name"`
	PlanningArea  string              `json:"planning_area"`
	StartLocation tripRequestLocation `json:"start_location"`
	StartDate     string              `json:"start_date"`
	EndDate       string              `json:"end_date"`
	Budget        tripRequestBudget   `json:"budget"`
	Interests     []string            `json:"interests"`
	TimeZone      string              `json:"time_zone"`
}

func newTripRequestBody(ownerID int, name string) tripRequestBody {
	startDate := time.Now().Add(48 * time.Hour)
	endDate := startDate.Add(48 * time.Hour)

	return tripRequestBody{
		OwnerID:      ownerID,
		Name:         name,
		PlanningArea: "Moscow",
		StartLocation: tripRequestLocation{
			Name:      "Red Square",
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

	body := newTripRequestBody(ownerID, name)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/trips/%d", baseURL, ownerID),
		marshalJSON(t, body),
	)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	response, err := newHTTPClient().Do(req)
	require.NoError(t, err)
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

func decodeErrorResponse(
	t *testing.T,
	response *http.Response,
) core_http_response.ErrorResponse {
	t.Helper()

	var result core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&result))
	return result
}

func newTripPatchBody(
	t *testing.T,
	ownerID int,
	version int,
	name string,
) *bytes.Reader {
	t.Helper()

	body := newTripRequestBody(ownerID, name)
	body.Version = version

	return marshalJSON(t, body)
}
