package users_transport_http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_errors "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/errors"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersTransportHTTP_GetUser_InvalidIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/users/invalid", nil)
	req.SetPathValue("userID", "invalid")

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_GetUser_NegativeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/users/-1", nil)
	req.SetPathValue("userID", "-1")

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_GetUser_ZeroID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/users/0", nil)
	req.SetPathValue("userID", "0")

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_GetUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	userID := 999
	expectedError := errors.New("database error")

	mockUsersService.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(nil, expectedError)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/users/"+strconv.Itoa(userID), nil)
	req.SetPathValue("userID", strconv.Itoa(userID))

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUsersTransportHTTP_GetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	userID := 999

	mockUsersService.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(nil, core_errors.ErrNotFound)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/users/"+strconv.Itoa(userID), nil)
	req.SetPathValue("userID", strconv.Itoa(userID))

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUsersTransportHTTP_GetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	userID := 1
	expectedUser := testutils.UserTest()

	mockUsersService.
		EXPECT().
		GetUser(gomock.Any(), userID).
		Return(&expectedUser, nil)

	rec, req := testutils.SetupHTTP(t, http.MethodGet, "/users/"+strconv.Itoa(userID), nil)
	req.SetPathValue("userID", strconv.Itoa(userID))

	handler.GetUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, testExpectedUserResponse(), response)
}
