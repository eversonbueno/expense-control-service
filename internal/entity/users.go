package entity

type User struct {
	ID        int     `json:"id"`
	Nome      string  `json:"nome"`
	Sobrenome string  `json:"sobrenome"`
	Usuario   string  `json:"usuario"`
	Senha     string  `json:"senha"`
	Saldo     float64 `json:"saldo"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
