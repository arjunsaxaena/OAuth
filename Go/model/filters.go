package model

import "github.com/gofrs/uuid"

type UserFilters struct {
	ID       *uuid.UUID `json:"id"`
	Name     *string    `json:"name"`
	Phone    *string    `json:"phone"`
	Email    *string    `json:"email"`
	Provider *string    `json:"provider"`
}