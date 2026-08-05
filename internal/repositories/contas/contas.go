package contas

import (
	"context"
	"database/sql"
	"errors"
	"expense-control-service/internal/entity"
	"fmt"
)

type Contas interface {
	Create(ctx context.Context, conta *entity.Conta) (*entity.Conta, error)
	FindByID(ctx context.Context, id int) (*entity.Conta, error)
	ListByGrupoFamiliar(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error)
	Update(ctx context.Context, conta *entity.Conta) (*entity.Conta, error)
	SetAtivo(ctx context.Context, id int, ativo bool) error
}

type contas struct {
	db *sql.DB
}

func New(db *sql.DB) Contas {
	return &contas{db: db}
}

func scanConta(row *sql.Row, conta *entity.Conta) error {
	return row.Scan(
		&conta.ID,
		&conta.GrupoFamiliarID,
		&conta.Nome,
		&conta.Tipo,
		&conta.FechamentoCartao,
		&conta.SaldoInicial,
		&conta.Ativo,
		&conta.CreatedAt,
		&conta.UpdatedAt,
	)
}

func (r *contas) Create(ctx context.Context, conta *entity.Conta) (*entity.Conta, error) {
	result, err := r.db.ExecContext(ctx, CreateConta,
		conta.GrupoFamiliarID, conta.Nome, conta.Tipo, conta.FechamentoCartao, conta.SaldoInicial,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar conta: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter id da conta criada: %v", err)
	}

	return r.FindByID(ctx, int(id))
}

// FindByID retorna nil, nil quando nenhuma conta é encontrada com o id informado.
func (r *contas) FindByID(ctx context.Context, id int) (*entity.Conta, error) {
	var conta entity.Conta
	row := r.db.QueryRowContext(ctx, FindContaById, id)

	if err := scanConta(row, &conta); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar conta: %v", err)
	}

	return &conta, nil
}

func (r *contas) ListByGrupoFamiliar(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error) {
	rows, err := r.db.QueryContext(ctx, ListContasByGrupoFamiliar, grupoFamiliarID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contas: %v", err)
	}
	defer rows.Close()

	var result []*entity.Conta
	for rows.Next() {
		var conta entity.Conta
		if err := rows.Scan(
			&conta.ID,
			&conta.GrupoFamiliarID,
			&conta.Nome,
			&conta.Tipo,
			&conta.FechamentoCartao,
			&conta.SaldoInicial,
			&conta.Ativo,
			&conta.CreatedAt,
			&conta.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao scanear conta: %v", err)
		}
		result = append(result, &conta)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após iterar contas: %v", err)
	}

	return result, nil
}

func (r *contas) Update(ctx context.Context, conta *entity.Conta) (*entity.Conta, error) {
	_, err := r.db.ExecContext(ctx, UpdateConta,
		conta.Nome, conta.Tipo, conta.FechamentoCartao, conta.SaldoInicial, conta.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar conta: %v", err)
	}

	return r.FindByID(ctx, conta.ID)
}

func (r *contas) SetAtivo(ctx context.Context, id int, ativo bool) error {
	_, err := r.db.ExecContext(ctx, SetContaAtivo, ativo, id)
	if err != nil {
		return fmt.Errorf("erro ao alterar status da conta: %v", err)
	}
	return nil
}
