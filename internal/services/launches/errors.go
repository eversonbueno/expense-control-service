package launches

import "errors"

var (
	// ErrTipoInvalido — RN-05: tipo diferente de receita/despesa.
	ErrTipoInvalido = errors.New("tipo do lançamento inválido")
	// ErrValorInvalido — RN-08: valor zero ou negativo.
	ErrValorInvalido = errors.New("valor deve ser maior que zero")
	// ErrContaInvalida — RN-04: conta inexistente, de outro grupo familiar ou desativada.
	ErrContaInvalida = errors.New("conta inexistente, de outro grupo familiar ou desativada")
	// ErrCategoriaInvalida — RN-06: categoria inexistente ou incompatível com o tipo informado.
	ErrCategoriaInvalida = errors.New("categoria inexistente ou incompatível com o tipo do lançamento")
	// ErrFormaPagamentoInvalida — RN-07: forma de pagamento fora da lista pré-definida.
	ErrFormaPagamentoInvalida = errors.New("forma de pagamento inválida")
	// ErrLancamentoNaoEncontrado — lançamento inexistente.
	ErrLancamentoNaoEncontrado = errors.New("lançamento não encontrado")
	// ErrAcessoNegado — CA-11: lançamento existe mas pertence a outro grupo familiar.
	ErrAcessoNegado = errors.New("lançamento pertence a outro grupo familiar")
)
