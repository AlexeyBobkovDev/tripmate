package testutils

import (
	"fmt"
	"strings"
	"time"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func BirthDateTest() time.Time {
	return time.Date(
		2010, 6, 9,
		0, 0, 0, 0,
		time.UTC,
	)
}

func CreatedAtTest() time.Time {
	return time.Date(
		2026, 9, 9,
		17, 0, 45, 0,
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
		Name:        TestName,
		Surname:     TestSurname,
		Username:    TestUsername,
		BirthDate:   birthDate,
		Description: TestDescription,
		Email:       TestEmail,
		PhoneNumber: TestPhoneNumber,
	}
}

func UserTest() domain.User {
	birthDate := BirthDateTest()
	return domain.User{
		ID:          1,
		Version:     1,
		Name:        TestName,
		Surname:     TestSurname,
		Username:    TestUsername,
		BirthDate:   birthDate,
		Description: TestDescription,
		Email:       TestEmail,
		PhoneNumber: TestPhoneNumber,
		CreatedAt:   CreatedAtTest(),
		UpdatedAt:   UpdatedAtTest(),
		DeletedAt:   nil,
	}
}

func UserPatchTest() domain.UserPatch {
	birthDate := BirthDateTest()

	return domain.UserPatch{
		Version:     1,
		Name:        domain.Nullable[string]{Value: &TestName, Set: true},
		Surname:     domain.Nullable[string]{Value: &TestSurname, Set: true},
		Username:    domain.Nullable[string]{Value: &TestUsername, Set: true},
		BirthDate:   domain.Nullable[time.Time]{Value: &birthDate, Set: true},
		Description: domain.Nullable[string]{Value: &TestDescription, Set: true},
		Email:       domain.Nullable[string]{Value: &TestEmail, Set: true},
		PhoneNumber: domain.Nullable[string]{Value: &TestPhoneNumber, Set: true},
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
		TestName,
		TestSurname,
		TestUsername,
		BirthDateTest().Format(core_consts.Layout),
		TestDescription,
		TestEmail,
		TestPhoneNumber,
	)

	return strings.NewReader(body)
}
