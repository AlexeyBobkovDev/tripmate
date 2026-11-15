package users_transport_http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	core_http_types "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/transport/http/types"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http/mocks"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func TestUsersTransportHTTP_PatchUser_InvalidIDFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/invalid", nil)
	req.SetPathValue("userID", "invalid")

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_NegativeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/-1", nil)
	req.SetPathValue("userID", "-1")

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_ZeroID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/0", nil)
	req.SetPathValue("userID", "0")

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/1", strings.NewReader("invalid json"))
	req.SetPathValue("userID", "1")

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_InvalidBirthDateFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	body := `{
		"version": 1,
		"birth_date": "invalid-date"
	}`
	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/1", strings.NewReader(body))
	req.SetPathValue("userID", "1")

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_BirthDateSetButNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	body := `{
		"version": 1,
		"birth_date": null
	}`
	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/1", strings.NewReader(body))
	req.SetPathValue("userID", "1")

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPatchUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request PatchUserRequest
		wantErr bool
	}{
		{
			name:    "valid request with no birth date",
			request: PatchUserRequest{Version: 1},
			wantErr: false,
		},
		{
			name:    "valid request with valid birth date",
			request: PatchUserRequest{Version: 1, BirthDate: core_http_types.Nullable[string]{Value: &[]string{"2006-01-02"}[0], Set: true}},
			wantErr: false,
		},
		{
			name:    "invalid birth date format",
			request: PatchUserRequest{Version: 1, BirthDate: core_http_types.Nullable[string]{Value: &[]string{"invalid-date"}[0], Set: true}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchUserRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersTransportHTTP_PatchUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	userID := 1
	body := `{
		"version": 1,
		"name": "Updated"
	}`
	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/"+strconv.Itoa(userID), strings.NewReader(body))
	req.SetPathValue("userID", strconv.Itoa(userID))

	expectedError := errors.New("service error")
	mockUsersService.
		EXPECT().
		PatchUser(gomock.Any(), userID, gomock.Any()).
		Return(nil, expectedError)

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	userID := 1
	body := `{
		"version": 1,
		"name": "Updated"
	}`
	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/"+strconv.Itoa(userID), strings.NewReader(body))
	req.SetPathValue("userID", strconv.Itoa(userID))

	conflictError := errors.New("conflict")
	mockUsersService.
		EXPECT().
		PatchUser(gomock.Any(), userID, gomock.Any()).
		Return(nil, conflictError)

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUsersTransportHTTP_PatchUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsersService := mocks.NewMockUsersService(ctrl)
	handler := NewUsersHTTPHandler(mockUsersService)

	userID := 1
	body := `{
		"version": 1,
		"name": "Updated",
		"email": "updated@example.com"
	}`
	rec, req := testutils.SetupHTTP(t, http.MethodPatch, "/users/"+strconv.Itoa(userID), strings.NewReader(body))
	req.SetPathValue("userID", strconv.Itoa(userID))

	expectedUser := testutils.UserTest()
	expectedUser.Name = "Updated"
	expectedUser.Email = "updated@example.com"

	mockUsersService.
		EXPECT().
		PatchUser(gomock.Any(), userID, gomock.Any()).
		Return(&expectedUser, nil)

	handler.PatchUser(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Name, response.Name)
	assert.Equal(t, expectedUser.Email, response.Email)
}
