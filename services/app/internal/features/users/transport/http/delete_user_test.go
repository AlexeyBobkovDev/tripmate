package users_transport_http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestDeleteUser_InvalidID(t *testing.T) {
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
				DeleteUser(
					gomock.Any(),
					gomock.Any(),
				).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodDelete, fmt.Sprintf("/api/v1/users/%s", tt.id), nil)
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("userID", tt.id)
			req = req.WithContext(ctx)

			usersTransport.DeleteUser(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		},
		)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.EXPECT().
		DeleteUser(
			gomock.Any(),
			1,
		).
		Return(
			nil,
		)

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/api/v1/users/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.DeleteUser(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestDeleteUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.EXPECT().
		DeleteUser(
			gomock.Any(),
			1,
		).
		Return(
			errors.New("service error"),
		)

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/api/v1/users/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.DeleteUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.EXPECT().
		DeleteUser(
			gomock.Any(),
			1,
		).
		Return(
			core_errors.ErrNotFound,
		)

	rec, req := testutils.SetupHTTP(t, http.MethodDelete, "/api/v1/users/1", nil)
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.DeleteUser(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
