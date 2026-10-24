package users_transport_http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersTransportHTTP_PatchUser_InvalidID(t *testing.T) {
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
				PatchUser(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodPatch, fmt.Sprintf("/api/v1/users/%s", tt.id), nil)
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)

			req.SetPathValue("userID", tt.id)
			req = req.WithContext(ctx)

			usersTransport.PatchUser(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
		},
		)
	}
}

func TestUsersTransportHTTP_PatchUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	userPatch := testutils.UserPatchTest()
	response := testUserResponse()
	expectedResponse := testExpectedUserResponse()

	mockUsersService.EXPECT().
		PatchUser(
			gomock.Any(),
			1,
			&userPatch,
		).
		Return(
			&response,
			nil,
		)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/api/v1/users/1", testPatchRequestBody())
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.PatchUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var actualUser UserResponse
	err := json.NewDecoder(rec.Body).Decode(&actualUser)
	require.NoError(t, err)

	assert.Equal(t, expectedResponse, actualUser)
}

func TestUsersTransportHTTP_PatchUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	request := testutils.UserPatchTest()
	serviceError := fmt.Errorf("service error")
	mockUsersService.EXPECT().
		PatchUser(
			gomock.Any(),
			1,
			&request,
		).
		Return(
			nil,
			serviceError,
		)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/api/v1/users/1", testPatchRequestBody())
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.PatchUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.EXPECT().
		PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/api/v1/users/1", testRequestInvalidJSON())
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)
	req.SetPathValue("userID", "1")

	usersTransport.PatchUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestUsersTransportHTTP_PatchUser_ValidationError(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "empty name",
			body: `{
				"name": "",
				"surname": "Bobkov",
				"username": "alex",
				"birth_date": "2010-06-09",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "+79998887766",
				"password": "password123"
			}`,
		},
		{
			name: "empty surname",
			body: `{
				"name": "Alex",
				"surname": "",
				"username": "alex",
				"birth_date": "2010-06-09",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "+79998887766",
				"password": "robust-password-Ab!(*&^1324"
			}`,
		},
		{
			name: "empty username",
			body: `{
				"name": "Alex",
				"surname": "Bobkov",
				"username": "",
				"birth_date": "2010-06-09",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "+79998887766",
				"password": "robust-password-Ab!(*&^1324"
			}`,
		},
		{
			name: "invalid email",
			body: `{
				"name": "Alex",
				"surname": "Bobkov",
				"username": "alex",
				"birth_date": "2010-06-09",
				"description": "description",
				"email": "alex",
				"phone_number": "+79998887766",
				"password": "robust-password-Ab!(*&^1324"
			}`,
		},
		{
			name: "invalid phone",
			body: `{
				"name": "Alex",
				"surname": "Bobkov",
				"username": "alex",
				"birth_date": "2010-06-09",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "123",
				"password": "robust-password-Ab!(*&^1324"
			}`,
		},

		{
			name: "invalid birth date format",
			body: `{
				"name": "Alex",
				"surname": "Bobkov",
				"username": "alex",
				"birth_date": "09-06-2010",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "+79998887766",
				"password": "robust-password-Ab!(*&^1324"
			}`,
		},
		{
			name: "future birth date",
			body: `{
				"name": "Alex",
				"surname": "Bobkov",
				"username": "alex",
				"birth_date": "2030-01-01",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "+79998887766",
				"password": "robust-password-Ab!(*&^1324"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsersService := mocks.NewMockUsersService(ctrl)
			usersTransport := NewUsersHTTPHandler(mockUsersService)

			mockUsersService.EXPECT().
				PatchUser(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/api/v1/users/1", strings.NewReader(tt.body))
			logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)
			req.SetPathValue("userID", "1")

			usersTransport.PatchUser(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}

func TestUsersTransportHTTP_PatchUser_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	request := testutils.UserPatchTest()
	conflictError := errors.Join(errors.New("conflict error"), core_errors.ErrConflict)
	mockUsersService.EXPECT().
		PatchUser(
			gomock.Any(),
			1,
			&request,
		).
		Return(
			nil,
			conflictError,
		)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/api/v1/users/1", testPatchRequestBody())
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.PatchUser(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var response core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&response))
	assert.Equal(t, conflictError.Error(), response.Error)
}

func TestUsersTransportHTTP_PatchUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	request := testutils.UserPatchTest()
	notFoundError := errors.Join(errors.New("not found"), core_errors.ErrNotFound)
	mockUsersService.EXPECT().
		PatchUser(
			gomock.Any(),
			1,
			&request,
		).
		Return(
			nil,
			notFoundError,
		)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/api/v1/users/1", testPatchRequestBody())
	logger := testutils.SetupLogger(t, testutils.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)

	req.SetPathValue("userID", "1")
	req = req.WithContext(ctx)

	usersTransport.PatchUser(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response core_http_response.ErrorResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&response))
	assert.Equal(t, notFoundError.Error(), response.Error)
}
