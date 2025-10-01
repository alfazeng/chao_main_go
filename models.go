package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	FullName     string    `json:"fullName"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // No exponer el hash en las respuestas JSON
	Country      string    `json:"country"`
	Phone        string    `json:"phone"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}