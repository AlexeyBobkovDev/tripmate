package users_transport_http

import "time"

type UserResponse struct {
	ID          int        `json:"id"             example:"1"`
	Version     int        `json:"version"        example:"5"`
	Name        string     `json:"name"           example:"Name"`
	Surname     string     `json:"surname"        example:"Surname"`
	Username    string     `json:"username"       example:"Username"`
	Description string     `json:"description"    example:"Some kind of description"`
	BirthDate   string     `json:"birth_date"     example:"2006-01-02"`
	Email       string     `json:"email"          example:"checkemail@gmail.com"`
	PhoneNumber string     `json:"phone_number"   example:"+79990746978"`
	CreatedAt   time.Time  `json:"created_at"     example:"2006-01-02T15:06:07.292454Z"`
	UpdatedAt   time.Time  `json:"updated_at"     example:"2006-01-02T15:06:07.292454Z"`
	DeletedAt   *time.Time `json:"deleted_at"     example:"2006-01-02T15:06:07.292454Z"`
}
