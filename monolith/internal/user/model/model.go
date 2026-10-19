package model

import (
	"time"
)

type User struct {
	ID           string     `db:"id"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required, min=2, max=50"`
	LastName  string `json:"last_name" validate:"required, min=2, max=50"`
	Email     string `json:"email"	validate:"required, email"`
	Password  string `json:"password"	validate:"required, min=8"`
	IsActive  bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
