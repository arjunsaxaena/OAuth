package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Provider string

const (
	PHONE Provider = "PHONE"
	GOOGLE Provider = "GOOGLE"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         *string    `json:"name"`
	Phone        *string    `json:"phone"`
	Email        *string    `json:"email"`
	Provider     Provider  `json:"provider"`
	ProfilePicture *string    `json:"profile_picture"`
	Meta         *string    `json:"meta"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}