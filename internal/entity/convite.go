package entity

import "time"

type Convite struct {
	ID              int
	Codigo          string
	GrupoFamiliarID int
	CriadoEm        time.Time
	ExpiraEm        time.Time
	UtilizadoEm     *time.Time
}

func (c *Convite) Utilizado() bool {
	return c.UtilizadoEm != nil
}

func (c *Convite) Expirado(agora time.Time) bool {
	return agora.After(c.ExpiraEm)
}
