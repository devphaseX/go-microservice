// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           pgtype.UUID        `json:"id"`
	Email        string             `json:"email"`
	FirstName    pgtype.Text        `json:"first_name"`
	LastName     pgtype.Text        `json:"last_name"`
	Password     string             `json:"password"`
	PasswordSalt string             `json:"password_salt"`
	UserActive   pgtype.Bool        `json:"user_active"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
}
