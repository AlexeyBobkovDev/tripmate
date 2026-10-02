package users_transport_http

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func mustParseDate(date string) time.Time {
	parsed, err := time.Parse(core_consts.Layout, date)
	if err != nil {
		panic(err)
	}
	return parsed
}

func TestUserDTOFromDomain(t *testing.T) {
	createdAt, updatedAt, deletedAt := time.Now(), time.Now(), time.Now()
	tests := []struct {
		name           string
		input          *domain.User
		expectedOutput UserResponse
	}{
		{
			name: "active user",
			input: &domain.User{
				ID:          1,
				Version:     5,
				Name:        "Name",
				Surname:     "Surname",
				Username:    "username",
				Description: "description",
				BirthDate:   mustParseDate("2006-01-02"),
				Email:       "some-email@gmail.com",
				PhoneNumber: "+79998887766",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				DeletedAt:   nil,
			},
			expectedOutput: UserResponse{
				ID:          1,
				Version:     5,
				Name:        "Name",
				Surname:     "Surname",
				Username:    "username",
				Description: "description",
				BirthDate:   "2006-01-02",
				Email:       "some-email@gmail.com",
				PhoneNumber: "+79998887766",
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				DeletedAt:   nil,
			},
		},
		{
			name: "deleted user",
			input: &domain.User{
				ID:        2,
				Version:   10,
				Name:      "Deleted",
				Surname:   "User",
				Username:  "deleted_user",
				BirthDate: mustParseDate("1999-12-31"),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				DeletedAt: &deletedAt,
			},
			expectedOutput: UserResponse{
				ID:        2,
				Version:   10,
				Name:      "Deleted",
				Surname:   "User",
				Username:  "deleted_user",
				BirthDate: "1999-12-31",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				DeletedAt: &deletedAt,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				output := userDTOFromDomain(tt.input)
				assert.Equal(t, tt.expectedOutput, output)
			},
		)
	}
}
