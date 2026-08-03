package requests

type RegistrarRequest struct {
	Nome          string  `json:"nome" binding:"required"`
	Email         string  `json:"email" binding:"required,email"`
	Senha         string  `json:"senha" binding:"required"`
	CodigoConvite *string `json:"codigo_convite"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required"`
}
