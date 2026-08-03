package auth

import "errors"

var (
	// ErrEmailJaCadastrado — RN-04: e-mail já cadastrado (handler mapeia para HTTP 409).
	ErrEmailJaCadastrado = errors.New("e-mail já cadastrado")

	// ErrSenhaCurta — RN-05: senha com menos de 8 caracteres (handler mapeia para HTTP 400).
	ErrSenhaCurta = errors.New("senha deve ter no mínimo 8 caracteres")

	// ErrConvite é a base para os erros de convite inválido (RN-10, handler mapeia para HTTP 400).
	// Use errors.Is(err, ErrConvite) para identificar a categoria; err.Error() traz o motivo específico.
	ErrConvite = errors.New("código de convite inválido")

	// ErrCredenciaisInvalidas — login com e-mail ou senha incorretos (handler mapeia para HTTP 401).
	ErrCredenciaisInvalidas = errors.New("e-mail ou senha inválidos")
)
