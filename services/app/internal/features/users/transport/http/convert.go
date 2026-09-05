package users_transport_http

import (
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

const layout = "2006-01-02"

func userDTOFromDomain(user *domain.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Version:     user.Version,
		Name:        user.Name,
		Surname:     user.Surname,
		Username:    user.Username,
		Description: user.Description,
		BirthDate:   user.BirthDate.Format(layout),
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
}
