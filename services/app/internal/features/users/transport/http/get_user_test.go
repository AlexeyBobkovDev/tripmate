package users_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testconsts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func TestGetUser_InvalidID(t *testing.T) {
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
				GetUser(
					gomock.Any(),
					gomock.Any(),
				).
				Times(0)

			rec, req := testhelpers.SetupHTTP(t, http.MethodGet, fmt.Sprintf("/api/v1/users/%s", tt.id), nil)
			logger := testhelpers.SetupLogger(t, testconsts.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("userID", tt.id)
			req = req.WithContext(ctx)

			usersTransport.GetUser(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		},
		)
	}
}

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	response := testUserResponse()
	mockUsersService.EXPECT().
		GetUser(
			gomock.Any(),
			1,
		).
		Return(
			&response,
			nil,
		)

	rec, req := testhelpers.SetupHTTP(t, http.MethodGet, "/api/v1/users/1", nil)
	logger := testhelpers.SetupLogger(t, testconsts.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.GetUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body UserResponse
	err := json.NewDecoder(rec.Body).Decode(&body)
	require.NoError(t, err)

	assert.Equal(t, testExpectedUserResponse(), body)
}

func TestGetUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.EXPECT().
		GetUser(
			gomock.Any(),
			1,
		).
		Return(
			nil,
			core_errors.ErrBadRequest,
		)

	rec, req := testhelpers.SetupHTTP(t, http.MethodGet, "/api/v1/users/1", nil)
	logger := testhelpers.SetupLogger(t, testconsts.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.GetUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
