package users_transport_http

import (
	"fmt"
	"strings"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testconsts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testhelpers"
)

func testUserResponse() domain.User {
	return domain.User{
		ID:          1,
		Version:     1,
		Name:        testconsts.TestName,
		Surname:     testconsts.TestSurname,
		Username:    testconsts.TestUsername,
		BirthDate:   testhelpers.BirthDateTest(),
		Description: testconsts.TestDescription,
		Email:       testconsts.TestEmail,
		PhoneNumber: testconsts.TestPhoneNumber,
		CreatedAt:   testhelpers.CreatedAtTest(),
		UpdatedAt:   testhelpers.UpdatedAtTest(),
		DeletedAt:   nil,
	}
}

func testExpectedUserResponse() UserResponse {
	expectedUser := testUserResponse()
	return UserResponse{
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
		testconsts.TestName,
		testconsts.TestSurname,
		testconsts.TestUsername,
		testhelpers.BirthDateTest().Format(core_consts.Layout),
		testconsts.TestDescription,
		testconsts.TestEmail,
		testconsts.TestPhoneNumber,
		testconsts.TestPassword,
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
		testconsts.TestName,
		testconsts.TestSurname,
		testconsts.TestUsername,
		testhelpers.BirthDateTest().Format(core_consts.Layout),
		testconsts.TestDescription,
		testconsts.TestEmail,
		testconsts.TestPhoneNumber,
	)

	return strings.NewReader(body)
}
