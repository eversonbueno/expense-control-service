package entity

import "time"

type GrupoFamiliar struct {
	ID       int       `json:"id"`
	CriadoEm time.Time `json:"criado_em"`
}
