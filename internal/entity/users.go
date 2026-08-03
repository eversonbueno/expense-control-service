package entity

import "time"

type User struct {
	ID              int       `json:"id"`
	Nome            string    `json:"nome"`
	Email           string    `json:"email"`
	SenhaHash       string    `json:"-"`
	GrupoFamiliarID int       `json:"grupo_familiar_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
