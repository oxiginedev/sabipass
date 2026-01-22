package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	ID              string     `bun:"type:uuid,pk" json:"id"`
	Name            *string    `bun:",nullzero" json:"name"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	EmailVerifiedAt *time.Time `bun:",nullzero" json:"email_verified_at"`
	Password        *string    `bun:",nullzero" json:"-"`
	Avatar          *string    `bun:",nullzero" json:"avatar"`
	GoogleID        *string    `bun:",nullzero" json:"-"`
	CreatedAt       time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`

	bun.BaseModel `bun:"table:users" json:"-"`
}

type FindUserOptions struct {
	ID    string
	Email string
}

type UserRepository interface {
	Create(context.Context, *User) error
	FindOne(context.Context, *FindUserOptions) (*User, error)
}
