//go:build e2e

package users_e2e_test

import (
	"fmt"
	"strings"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	users_transport_http "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/users/transport/http"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testutils"
)

func testUserResponse() domain.User {
	return domain.User{
		ID:          1,
		Version:     1,
		Name:        testutils.TestName,
		Surname:     testutils.TestSurname,
		Username:    testutils.TestUsername,
		BirthDate:   testutils.BirthDateTest(),
		Description: testutils.TestDescription,
		Email:       testutils.TestEmail,
		PhoneNumber: testutils.TestPhoneNumber,
		CreatedAt:   testutils.CreatedAtTest(),
		UpdatedAt:   testutils.UpdatedAtTest(),
		DeletedAt:   nil,
	}
}

func testExpectedUserResponse() users_transport_http.UserResponse {
	expectedUser := testUserResponse()
	return users_transport_http.UserResponse{
		ID:          expectedUser.ID,
		Version:     expectedUser.Version,
		Name:        expectedUser.Name,
		Surname:     expectedUser.Surname,
		Username:    expectedUser.Username,
		Description: expectedUser.Description,
		BirthDate:   expectedUser.BirthDate.Format(core_consts.Layout),
		Email:       expectedUser.Email,
		PhoneNumber: expectedUser.PhoneNumber,
		CreatedAt:   expectedUser.CreatedAt,
		UpdatedAt:   expectedUser.UpdatedAt,
		DeletedAt:   expectedUser.DeletedAt,
	}
}

func testRequestBodyWithDefaults() *strings.Reader {
	body := fmt.Sprintf(`
		{
			"name": "%s",
			"surname": "%s",
			"username": "%s",
			"birth_date": "%s",
			"description": "%s",
			"email": "%s",
			"phone_number": "%s",
			"password": "%s"
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
	return strings.NewReader(body)
}

func testRequestInvalidJSON() *strings.Reader {
	return strings.NewReader(`{
			"name": "Name"
			"surname": "Surname"
			"username": "username"
			"birth_date": "2010-06-09"
			"description": "description"
			"email": "some-email@gmail.com"
			"phone_number": "+79998887766"
			"password": "some-password"
		}`)
}

func testPatchRequestBody() *strings.Reader {
	body := fmt.Sprintf(`
		{
			"version": 1,
			"name": "%s",
			"surname": "%s",
			"username": "%s",
			"birth_date": "%s",
			"description": "%s",
			"email": "%s",
			"phone_number": "%s"
		}`,
		testutils.TestName,
		testutils.TestSurname,
		testutils.TestUsername,
		testutils.BirthDateTest().Format(core_consts.Layout),
		testutils.TestDescription,
		testutils.TestEmail,
		testutils.TestPhoneNumber,
	)

	return strings.NewReader(body)
}

func testPatchRequestBodyWithIncreasedByOneVersion(version int) *strings.Reader {
	body := fmt.Sprintf(`
		{
			"version": %d,
			"name": "%s",
			"surname": "%s",
			"username": "%s",
			"birth_date": "%s",
			"description": "%s",
			"email": "%s",
			"phone_number": "%s"
		}`,
		version+1,
		testutils.TestName,
		testutils.TestSurname,
		testutils.TestUsername,
		testutils.BirthDateTest().Format(core_consts.Layout),
		testutils.TestDescription,
		testutils.TestEmail,
		testutils.TestPhoneNumber,
	)

	return strings.NewReader(body)
}
