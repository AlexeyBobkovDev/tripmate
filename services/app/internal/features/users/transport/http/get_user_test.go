package users_transport_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersTransportHTTP_GetUser_InvalidID(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{
			name: "empty ID",
			id:   "",
		},
		{
			name: "letter in ID",
			id:   "asdf",
		},
		{
			name: "special symbols in ID",
			id:   "&^$*^&",
		},
		{
			name: "negative ID",
			id:   "-1",
		},
		{
			name: "zero ID",
			id:   "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsersService := mocks.NewMockUsersService(ctrl)
			usersTransport := NewUsersHTTPHandler(mockUsersService)

			mockUsersService.EXPECT().
				GetUser(gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(
				t,
				http.MethodGet,
				fmt.Sprintf("/api/v1/users/%s", tt.id),
				nil,
			)

			logger := testutils.SetupLogger(
				t,
				testutils.TestLevel,
				t.TempDir(),
			)
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("userID", tt.id)
			req = req.WithContext(ctx)

			usersTransport.GetUser(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestUsersTransportHTTP_GetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	expectedUser := testUserResponse()

	mockUsersService.EXPECT().
		GetUser(gomock.Any(), 1).
		Return(&expectedUser, nil)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodGet,
		"/api/v1/users/1",
		nil,
	)

	logger := testutils.SetupLogger(
		t,
		testutils.TestLevel,
		t.TempDir(),
	)
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.GetUser(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body UserResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

	assert.Equal(t, testExpectedUserResponse(), body)
}

func TestUsersTransportHTTP_GetUser_BadRequestServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.EXPECT().
		GetUser(gomock.Any(), 1).
		Return(nil, core_errors.ErrBadRequest)

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodGet,
		"/api/v1/users/1",
		nil,
	)

	logger := testutils.SetupLogger(
		t,
		testutils.TestLevel,
		t.TempDir(),
	)
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.GetUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_GetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.EXPECT().
		GetUser(gomock.Any(), 1).
		Return(nil, fmt.Errorf("get userID %d: %w", 1, errors.Join(pgx.ErrNoRows, core_errors.ErrNotFound)))

	rec, req := testutils.SetupHTTP(
		t,
		http.MethodGet,
		"/api/v1/users/1",
		nil,
	)

	logger := testutils.SetupLogger(
		t,
		testutils.TestLevel,
		t.TempDir(),
	)
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.GetUser(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
