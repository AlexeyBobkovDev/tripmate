package testhelpers

import (
	"fmt"
	"strings"
	"time"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/testconsts"
)

func BirthDateTest() time.Time {
	return time.Date(
		2010, 0o6, 9,
		0, 0, 0, 0,
		time.UTC,
	)
}

func CreatedAtTest() time.Time {
	return time.Date(
		2026, 9, 9,
		17, 0o0, 45, 0,
		time.UTC,
	)
}

func UpdatedAtTest() time.Time {
	return time.Date(
		2026, 9, 9,
		17, 10, 45, 0,
		time.UTC,
	)
}

func UserCreateTest() domain.UserCreate {
	birthDate := BirthDateTest()
	return domain.UserCreate{
		Name:        testconsts.TestName,
		Surname:     testconsts.TestSurname,
		Username:    testconsts.TestUsername,
		BirthDate:   birthDate,
		Description: testconsts.TestDescription,
		Email:       testconsts.TestEmail,
		PhoneNumber: testconsts.TestPhoneNumber,
	}
}

func UserPatchTest() domain.UserPatch {
	birthDate := BirthDateTest()

	return domain.UserPatch{
		Version:     1,
		Name:        domain.Nullable[string]{Value: &testconsts.TestName, Set: true},
		Surname:     domain.Nullable[string]{Value: &testconsts.TestSurname, Set: true},
		Username:    domain.Nullable[string]{Value: &testconsts.TestUsername, Set: true},
		BirthDate:   domain.Nullable[time.Time]{Value: &birthDate, Set: true},
		Description: domain.Nullable[string]{Value: &testconsts.TestDescription, Set: true},
		Email:       domain.Nullable[string]{Value: &testconsts.TestEmail, Set: true},
		PhoneNumber: domain.Nullable[string]{Value: &testconsts.TestPhoneNumber, Set: true},
	}
}

func PatchRequestBodyTest() *strings.Reader {
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
		BirthDateTest().Format(core_consts.Layout),
		testconsts.TestDescription,
		testconsts.TestEmail,
		testconsts.TestPhoneNumber,
	)

	return strings.NewReader(body)
}
