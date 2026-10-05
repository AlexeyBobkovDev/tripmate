package users_transport_http

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_logger "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/logger"
	core_postgres_pool_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/repository/postgres/pool/errors"
	core_http_response "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/response"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testconsts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func TestCreateUser_Success(t *testing.T) {
	userCreate := testhelpers.UserCreateTest()
	userResponse := testUserResponse()

	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.
		EXPECT().
		CreateUser(
			gomock.Any(),
			&userCreate,
			testconsts.TestPassword,
		).
		Return(&userResponse, nil)

	rec, req := testhelpers.SetupHTTP(
		t,
		http.MethodPost,
		"/api/v1/users",
		testRequestBodyWithDefaults(),
	)

	logger := testhelpers.SetupLogger(t, testconsts.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	usersTransport.CreateUser(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var user UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testExpectedUserResponse(), user)
}

func TestCreateUser_ServiceError(t *testing.T) {
	serviceError := core_postgres_pool_errors.ErrViolatedForeignKey
	userCreate := testhelpers.UserCreateTest()

	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.
		EXPECT().
		CreateUser(
			gomock.Any(),
			&userCreate,
			testconsts.TestPassword,
		).
		Return(nil, serviceError)

	rec, req := testhelpers.SetupHTTP(
		t,
		http.MethodPost,
		"/api/v1/users",
		testRequestBodyWithDefaults(),
	)

	logger := testhelpers.SetupLogger(t, testconsts.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	usersTransport.CreateUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var returnedServiceError core_http_response.ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&returnedServiceError); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "failed to create user", returnedServiceError.Message)
	assert.Equal(t, core_postgres_pool_errors.ErrViolatedForeignKey.Error(), returnedServiceError.Error)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUsersService := mocks.NewMockUsersService(ctrl)
	usersTransport := NewUsersHTTPHandler(mockUsersService)

	mockUsersService.
		EXPECT().
		CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	rec, req := testhelpers.SetupHTTP(
		t,
		http.MethodPost,
		"/api/v1/users",
		testRequestInvalidJSON(),
	)

	logger := testhelpers.SetupLogger(t, testconsts.TestLevel, t.TempDir())
	ctx := core_logger.ToContext(req.Context(), logger)
	req = req.WithContext(ctx)

	usersTransport.CreateUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func TestCreateUser_ValidationError(t *testing.T) {
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
				"password": "robust-password-Ab!(*&^1324"
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
			name: "password too short",
			body: `{
				"name": "Alex",
				"surname": "Bobkov",
				"username": "alex",
				"birth_date": "2010-06-09",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "+79998887766",
				"password": "123"
			}`,
		},
		{
			name: "empty password",
			body: `{
				"name": "Alex",
				"surname": "Bobkov",
				"username": "alex",
				"birth_date": "2010-06-09",
				"description": "description",
				"email": "alex@gmail.com",
				"phone_number": "+79998887766",
				"password": ""
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUsersService := mocks.NewMockUsersService(ctrl)
			usersTransport := NewUsersHTTPHandler(mockUsersService)

			mockUsersService.
				EXPECT().
				CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(0)

			rec, req := testhelpers.SetupHTTP(
				t,
				http.MethodPost,
				"/api/v1/users",
				strings.NewReader(tt.body),
			)

			logger := testhelpers.SetupLogger(
				t,
				testconsts.TestLevel,
				t.TempDir(),
			)

			ctx := core_logger.ToContext(req.Context(), logger)
			req = req.WithContext(ctx)

			usersTransport.CreateUser(rec, req)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		})
	}
}
