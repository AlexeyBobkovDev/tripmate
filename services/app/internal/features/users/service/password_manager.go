package users_service

//go:generate mockgen -source=password_manager.go -destination=mocks/user_password_manager.go -package mocks
type PasswordManager interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}
