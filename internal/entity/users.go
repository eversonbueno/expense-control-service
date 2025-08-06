package entity

type User struct {
	ID           int     `json:"id"`
	NomeCompleto string  `json:"nome_completo"`
	CPFCNPJ      string  `json:"cpf_cnpj"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	TipoUsuario  int     `json:"tipo_usuario"`
	SaldoUsuario float64 `json:"saldo_usuario"`
}
