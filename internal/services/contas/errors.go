package contas

import "errors"

var (
	// ErrTipoInvalido — RN-04: tipo fora da lista pré-definida.
	ErrTipoInvalido = errors.New("tipo de conta inválido")
	// ErrFechamentoCartaoInvalido — RN-05: fechamento_cartao informado para tipo diferente de cartao_credito.
	ErrFechamentoCartaoInvalido = errors.New("fechamento_cartao só é permitido para contas do tipo cartao_credito")
	// ErrContaNaoEncontrada — conta inexistente.
	ErrContaNaoEncontrada = errors.New("conta não encontrada")
	// ErrAcessoNegado — CA-08: conta existe mas pertence a outro grupo familiar.
	ErrAcessoNegado = errors.New("conta pertence a outro grupo familiar")
)
