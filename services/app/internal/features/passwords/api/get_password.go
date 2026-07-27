package passwords_api

import (
	"context"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/domain"
)

func (api *PasswordAPI) GetPassword(
	ctx context.Context,
	userID int,
) (*domain.Password, error) {
	return api.passwordService.GetPassword(ctx, userID)
}
