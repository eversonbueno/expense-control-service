package contas

import (
	"context"
	"expense-control-service/internal/entity"
	contasRepo "expense-control-service/internal/repositories/contas"
)

const (
	TipoCorrente      = "corrente"
	TipoPoupanca      = "poupanca"
	TipoCartaoCredito = "cartao_credito"
	TipoDinheiro      = "dinheiro"
)

var tiposValidos = map[string]bool{
	TipoCorrente:      true,
	TipoPoupanca:      true,
	TipoCartaoCredito: true,
	TipoDinheiro:      true,
}

type Contas interface {
	Criar(ctx context.Context, grupoFamiliarID int, nome, tipo string, fechamentoCartao *int, saldoInicial *float64) (*entity.Conta, error)
	Listar(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error)
	Atualizar(ctx context.Context, grupoFamiliarID, id int, nome, tipo string, fechamentoCartao *int, saldoInicial float64) (*entity.Conta, error)
	Desativar(ctx context.Context, grupoFamiliarID, id int) (*entity.Conta, error)
	Reativar(ctx context.Context, grupoFamiliarID, id int) (*entity.Conta, error)
}

type contas struct {
	repo contasRepo.Contas
}

func New(repo contasRepo.Contas) Contas {
	return &contas{repo: repo}
}

// validarTipoEFechamento aplica RN-04 (tipo restrito à lista pré-definida) e
// RN-05 (fechamento_cartao só aceito em contas do tipo cartao_credito).
func validarTipoEFechamento(tipo string, fechamentoCartao *int) error {
	if !tiposValidos[tipo] {
		return ErrTipoInvalido
	}
	if fechamentoCartao != nil && tipo != TipoCartaoCredito {
		return ErrFechamentoCartaoInvalido
	}
	return nil
}

func (s *contas) Criar(ctx context.Context, grupoFamiliarID int, nome, tipo string, fechamentoCartao *int, saldoInicial *float64) (*entity.Conta, error) {
	if err := validarTipoEFechamento(tipo, fechamentoCartao); err != nil {
		return nil, err
	}

	saldo := 0.0
	if saldoInicial != nil {
		saldo = *saldoInicial
	}

	conta := &entity.Conta{
		GrupoFamiliarID:  grupoFamiliarID,
		Nome:             nome,
		Tipo:             tipo,
		FechamentoCartao: fechamentoCartao,
		SaldoInicial:     saldo,
		Ativo:            true,
	}

	return s.repo.Create(ctx, conta)
}

// Listar retorna todas as contas do grupo familiar (ativas e inativas — CA-04).
func (s *contas) Listar(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error) {
	return s.repo.ListByGrupoFamiliar(ctx, grupoFamiliarID)
}

// buscarDoGrupo aplica CA-08: distingue conta inexistente (404) de conta
// existente porém de outro grupo familiar (403).
func (s *contas) buscarDoGrupo(ctx context.Context, grupoFamiliarID, id int) (*entity.Conta, error) {
	conta, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if conta == nil {
		return nil, ErrContaNaoEncontrada
	}
	if conta.GrupoFamiliarID != grupoFamiliarID {
		return nil, ErrAcessoNegado
	}
	return conta, nil
}

func (s *contas) Atualizar(ctx context.Context, grupoFamiliarID, id int, nome, tipo string, fechamentoCartao *int, saldoInicial float64) (*entity.Conta, error) {
	conta, err := s.buscarDoGrupo(ctx, grupoFamiliarID, id)
	if err != nil {
		return nil, err
	}

	if err := validarTipoEFechamento(tipo, fechamentoCartao); err != nil {
		return nil, err
	}

	conta.Nome = nome
	conta.Tipo = tipo
	conta.FechamentoCartao = fechamentoCartao
	conta.SaldoInicial = saldoInicial

	return s.repo.Update(ctx, conta)
}

func (s *contas) Desativar(ctx context.Context, grupoFamiliarID, id int) (*entity.Conta, error) {
	return s.setAtivo(ctx, grupoFamiliarID, id, false)
}

func (s *contas) Reativar(ctx context.Context, grupoFamiliarID, id int) (*entity.Conta, error) {
	return s.setAtivo(ctx, grupoFamiliarID, id, true)
}

func (s *contas) setAtivo(ctx context.Context, grupoFamiliarID, id int, ativo bool) (*entity.Conta, error) {
	conta, err := s.buscarDoGrupo(ctx, grupoFamiliarID, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SetAtivo(ctx, conta.ID, ativo); err != nil {
		return nil, err
	}

	conta.Ativo = ativo
	return conta, nil
}
