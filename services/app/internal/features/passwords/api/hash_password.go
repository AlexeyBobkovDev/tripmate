package passwords_api

import passwords_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/service"

func (api *PasswordAPI) HashPassword(
	password []byte,
	times uint32,
	memory uint32,
	threads uint8,
	keyLen uint32,
) ([]byte, []byte, error) {
	passwordHasher := passwords_service.NewArgon2IDHasher(
		times,
		memory,
		threads,
		keyLen,
	)
	passwordHash, salt, err := passwordHasher.Hash(password)
	return passwordHash, salt, err
}
