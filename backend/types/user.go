package types

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"-"` // Omit password from JSON responses
	PasswordHash string    `json:"-"` // Omit hash from JSON responses
	CreatedAt    time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
