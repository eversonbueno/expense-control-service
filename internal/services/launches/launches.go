package launches

import (
	"context"
	"time"

	"expense-control-service/internal/entity"
	contasRepo "expense-control-service/internal/repositories/contas"
	launchCategoryRepo "expense-control-service/internal/repositories/launch_category"
	launchesRepo "expense-control-service/internal/repositories/launches"
	paymentMethodsRepo "expense-control-service/internal/repositories/payment_methods"
)

const (
	TipoReceita = "receita"
	TipoDespesa = "despesa"
)

type Launches interface {
	ListLaunches(ctx context.Context, grupoFamiliarID int, mes, ano *int) ([]*entity.Launches, error)
	Criar(ctx context.Context, grupoFamiliarID, usuarioID, contaID int, tipo string, categoriaID, formaPagamentoID int, descricao string, valor float64, data time.Time) (*entity.Launches, error)
	Atualizar(ctx context.Context, grupoFamiliarID, id, contaID int, tipo string, categoriaID, formaPagamentoID int, descricao string, valor float64, data time.Time) (*entity.Launches, error)
	Excluir(ctx context.Context, grupoFamiliarID, id int) error
}

type launches struct {
	launchesRepo       launchesRepo.Launches
	contasRepo         contasRepo.Contas
	launchCategoryRepo launchCategoryRepo.LaunchCategory
	paymentMethodsRepo paymentMethodsRepo.PaymentMethods
}

func New(
	launchesRepo launchesRepo.Launches,
	contasRepo contasRepo.Contas,
	launchCategoryRepo launchCategoryRepo.LaunchCategory,
	paymentMethodsRepo paymentMethodsRepo.PaymentMethods,
) Launches {
	return &launches{
		launchesRepo:       launchesRepo,
		contasRepo:         contasRepo,
		launchCategoryRepo: launchCategoryRepo,
		paymentMethodsRepo: paymentMethodsRepo,
	}
}

// tipoLancamentoID aplica RN-05: tipo deve ser um dos dois valores fixos, mapeados para os ids
// já cadastrados em tipo_lancamento (1=Entrada/receita, 2=Saida/despesa).
func tipoLancamentoID(tipo string) (int, error) {
	switch tipo {
	case TipoReceita:
		return int(entity.LauchTypeEntrada.ID), nil
	case TipoDespesa:
		return int(entity.LauchTypeSaida.ID), nil
	default:
		return 0, ErrTipoInvalido
	}
}

func (l *launches) ListLaunches(ctx context.Context, grupoFamiliarID int, mes, ano *int) ([]*entity.Launches, error) {
	return l.launchesRepo.ListLaunches(ctx, grupoFamiliarID, mes, ano)
}

func (l *launches) Criar(ctx context.Context, grupoFamiliarID, usuarioID, contaID int, tipo string, categoriaID, formaPagamentoID int, descricao string, valor float64, data time.Time) (*entity.Launches, error) {
	tipoID, err := l.validar(ctx, grupoFamiliarID, contaID, tipo, categoriaID, formaPagamentoID, valor)
	if err != nil {
		return nil, err
	}

	launch := &entity.Launches{
		Usuario:             usuarioID,
		GrupoFamiliarID:     grupoFamiliarID,
		ContaID:             contaID,
		FormaPagamento:      formaPagamentoID,
		TipoLancamento:      tipoID,
		CategoriaLancamento: categoriaID,
		Data:                data,
		Descricao:           descricao,
		Valor:               valor,
	}

	return l.launchesRepo.Create(ctx, launch)
}

func (l *launches) Atualizar(ctx context.Context, grupoFamiliarID, id, contaID int, tipo string, categoriaID, formaPagamentoID int, descricao string, valor float64, data time.Time) (*entity.Launches, error) {
	launch, err := l.buscarDoGrupo(ctx, grupoFamiliarID, id)
	if err != nil {
		return nil, err
	}

	tipoID, err := l.validar(ctx, grupoFamiliarID, contaID, tipo, categoriaID, formaPagamentoID, valor)
	if err != nil {
		return nil, err
	}

	launch.ContaID = contaID
	launch.TipoLancamento = tipoID
	launch.CategoriaLancamento = categoriaID
	launch.FormaPagamento = formaPagamentoID
	launch.Descricao = descricao
	launch.Valor = valor
	launch.Data = data

	return l.launchesRepo.Update(ctx, launch)
}

func (l *launches) Excluir(ctx context.Context, grupoFamiliarID, id int) error {
	if _, err := l.buscarDoGrupo(ctx, grupoFamiliarID, id); err != nil {
		return err
	}
	return l.launchesRepo.SoftDelete(ctx, id)
}

// buscarDoGrupo aplica CA-11: distingue lançamento inexistente (404) de lançamento
// existente porém de outro grupo familiar (403) — mesmo padrão de services/contas.
func (l *launches) buscarDoGrupo(ctx context.Context, grupoFamiliarID, id int) (*entity.Launches, error) {
	launch, err := l.launchesRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if launch == nil {
		return nil, ErrLancamentoNaoEncontrado
	}
	if launch.GrupoFamiliarID != grupoFamiliarID {
		return nil, ErrAcessoNegado
	}
	return launch, nil
}

// validar aplica RN-05 a RN-08, comuns ao cadastro e à edição (RN-10).
func (l *launches) validar(ctx context.Context, grupoFamiliarID, contaID int, tipo string, categoriaID, formaPagamentoID int, valor float64) (int, error) {
	tipoID, err := tipoLancamentoID(tipo)
	if err != nil {
		return 0, err
	}

	if valor <= 0 {
		return 0, ErrValorInvalido
	}

	if err := l.validarConta(ctx, grupoFamiliarID, contaID); err != nil {
		return 0, err
	}

	if err := l.validarCategoria(ctx, categoriaID, tipoID); err != nil {
		return 0, err
	}

	if err := l.validarFormaPagamento(ctx, formaPagamentoID); err != nil {
		return 0, err
	}

	return tipoID, nil
}

// validarConta aplica RN-04: a conta deve existir, pertencer ao grupo familiar e estar ativa.
func (l *launches) validarConta(ctx context.Context, grupoFamiliarID, contaID int) error {
	conta, err := l.contasRepo.FindByID(ctx, contaID)
	if err != nil {
		return err
	}
	if conta == nil || conta.GrupoFamiliarID != grupoFamiliarID || !conta.Ativo {
		return ErrContaInvalida
	}
	return nil
}

// validarCategoria aplica RN-06: a categoria deve existir e ser compatível com o tipo informado.
func (l *launches) validarCategoria(ctx context.Context, categoriaID, tipoID int) error {
	categoria, err := l.launchCategoryRepo.FindByID(ctx, categoriaID)
	if err != nil {
		return err
	}
	if categoria == nil || int(categoria.TipoLancamentoID) != tipoID {
		return ErrCategoriaInvalida
	}
	return nil
}

// validarFormaPagamento aplica RN-07: a forma de pagamento deve existir na lista pré-definida.
func (l *launches) validarFormaPagamento(ctx context.Context, formaPagamentoID int) error {
	existe, err := l.paymentMethodsRepo.Exists(ctx, formaPagamentoID)
	if err != nil {
		return err
	}
	if !existe {
		return ErrFormaPagamentoInvalida
	}
	return nil
}
