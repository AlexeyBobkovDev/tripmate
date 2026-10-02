package users_transport_http

import (
	core_consts "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/consts"
	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func userDTOFromDomain(user *domain.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Version:     user.Version,
		Name:        user.Name,
		Surname:     user.Surname,
		Username:    user.Username,
		Description: user.Description,
		BirthDate:   user.BirthDate.Format(core_consts.Layout),
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
}
