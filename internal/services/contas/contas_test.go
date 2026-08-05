package contas

import (
	"context"
	"errors"
	"expense-control-service/internal/entity"
	"testing"
)

type mockContasRepo struct {
	createFn              func(ctx context.Context, conta *entity.Conta) (*entity.Conta, error)
	findByIDFn            func(ctx context.Context, id int) (*entity.Conta, error)
	listByGrupoFamiliarFn func(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error)
	updateFn              func(ctx context.Context, conta *entity.Conta) (*entity.Conta, error)
	setAtivoFn            func(ctx context.Context, id int, ativo bool) error
}

func (m *mockContasRepo) Create(ctx context.Context, conta *entity.Conta) (*entity.Conta, error) {
	return m.createFn(ctx, conta)
}
func (m *mockContasRepo) FindByID(ctx context.Context, id int) (*entity.Conta, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockContasRepo) ListByGrupoFamiliar(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error) {
	return m.listByGrupoFamiliarFn(ctx, grupoFamiliarID)
}
func (m *mockContasRepo) Update(ctx context.Context, conta *entity.Conta) (*entity.Conta, error) {
	return m.updateFn(ctx, conta)
}
func (m *mockContasRepo) SetAtivo(ctx context.Context, id int, ativo bool) error {
	return m.setAtivoFn(ctx, id, ativo)
}

func floatPtr(v float64) *float64 { return &v }
func intPtr(v int) *int           { return &v }

func TestCriar(t *testing.T) {
	tests := []struct {
		name             string
		grupoFamiliarID  int
		nomeConta        string
		tipo             string
		fechamentoCartao *int
		saldoInicial     *float64
		wantErr          error
		wantSaldo        float64
	}{
		{
			// CA-01
			name:            "conta valida e criada com sucesso",
			grupoFamiliarID: 1,
			nomeConta:       "Nubank",
			tipo:            TipoCartaoCredito,
			wantSaldo:       0,
		},
		{
			// CA-02
			name:            "saldo_inicial omitido assume zero",
			grupoFamiliarID: 1,
			nomeConta:       "Carteira",
			tipo:            TipoDinheiro,
			saldoInicial:    nil,
			wantSaldo:       0,
		},
		{
			name:            "saldo_inicial informado é respeitado",
			grupoFamiliarID: 1,
			nomeConta:       "Bradesco",
			tipo:            TipoCorrente,
			saldoInicial:    floatPtr(150.50),
			wantSaldo:       150.50,
		},
		{
			// CA-03
			name:            "tipo fora da lista pré-definida retorna ErrTipoInvalido",
			grupoFamiliarID: 1,
			nomeConta:       "X",
			tipo:            "invalido",
			wantErr:         ErrTipoInvalido,
		},
		{
			// CA-03b
			name:             "fechamento_cartao em tipo diferente de cartao_credito retorna ErrFechamentoCartaoInvalido",
			grupoFamiliarID:  1,
			nomeConta:        "X",
			tipo:             TipoCorrente,
			fechamentoCartao: intPtr(10),
			wantErr:          ErrFechamentoCartaoInvalido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockContasRepo{
				createFn: func(ctx context.Context, conta *entity.Conta) (*entity.Conta, error) {
					conta.ID = 1
					return conta, nil
				},
			}
			svc := New(repo)

			conta, err := svc.Criar(context.Background(), tt.grupoFamiliarID, tt.nomeConta, tt.tipo, tt.fechamentoCartao, tt.saldoInicial)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if conta.GrupoFamiliarID != tt.grupoFamiliarID {
				t.Fatalf("esperava grupo_familiar_id=%d, recebeu %d", tt.grupoFamiliarID, conta.GrupoFamiliarID)
			}
			if conta.SaldoInicial != tt.wantSaldo {
				t.Fatalf("esperava saldo_inicial=%v, recebeu %v", tt.wantSaldo, conta.SaldoInicial)
			}
			if !conta.Ativo {
				t.Fatal("esperava conta criada como ativa")
			}
		})
	}
}

func TestListar(t *testing.T) {
	// CA-04
	repo := &mockContasRepo{
		listByGrupoFamiliarFn: func(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error) {
			if grupoFamiliarID != 1 {
				t.Fatalf("esperava grupo_familiar_id=1, recebeu %d", grupoFamiliarID)
			}
			return []*entity.Conta{
				{ID: 1, GrupoFamiliarID: 1, Nome: "Nubank"},
				{ID: 2, GrupoFamiliarID: 1, Nome: "Carteira"},
			}, nil
		},
	}
	svc := New(repo)

	contas, err := svc.Listar(context.Background(), 1)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(contas) != 2 {
		t.Fatalf("esperava 2 contas, recebeu %d", len(contas))
	}
}

func TestAtualizar(t *testing.T) {
	contaExistente := &entity.Conta{ID: 1, GrupoFamiliarID: 1, Nome: "Nubank", Tipo: TipoCorrente, Ativo: true}

	tests := []struct {
		name             string
		grupoFamiliarID  int
		id               int
		tipo             string
		fechamentoCartao *int
		repo             *mockContasRepo
		wantErr          error
	}{
		{
			// CA-05
			name:            "edita conta do proprio grupo com sucesso",
			grupoFamiliarID: 1,
			id:              1,
			tipo:            TipoPoupanca,
			repo: &mockContasRepo{
				findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return contaExistente, nil },
				updateFn:   func(ctx context.Context, conta *entity.Conta) (*entity.Conta, error) { return conta, nil },
			},
		},
		{
			// CA-08: conta inexistente
			name:            "conta inexistente retorna ErrContaNaoEncontrada",
			grupoFamiliarID: 1,
			id:              999,
			tipo:            TipoCorrente,
			repo: &mockContasRepo{
				findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return nil, nil },
			},
			wantErr: ErrContaNaoEncontrada,
		},
		{
			// CA-08: conta de outro grupo familiar
			name:            "conta de outro grupo familiar retorna ErrAcessoNegado",
			grupoFamiliarID: 2,
			id:              1,
			tipo:            TipoCorrente,
			repo: &mockContasRepo{
				findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return contaExistente, nil },
			},
			wantErr: ErrAcessoNegado,
		},
		{
			// CA-03
			name:            "tipo invalido retorna ErrTipoInvalido mesmo apos encontrar a conta",
			grupoFamiliarID: 1,
			id:              1,
			tipo:            "invalido",
			repo: &mockContasRepo{
				findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return contaExistente, nil },
			},
			wantErr: ErrTipoInvalido,
		},
		{
			// CA-03b
			name:             "fechamento_cartao invalido na edicao retorna ErrFechamentoCartaoInvalido",
			grupoFamiliarID:  1,
			id:               1,
			tipo:             TipoCorrente,
			fechamentoCartao: intPtr(15),
			repo: &mockContasRepo{
				findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return contaExistente, nil },
			},
			wantErr: ErrFechamentoCartaoInvalido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.repo)

			_, err := svc.Atualizar(context.Background(), tt.grupoFamiliarID, tt.id, "Novo nome", tt.tipo, tt.fechamentoCartao, 100)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

func TestDesativarReativar(t *testing.T) {
	contaExistente := &entity.Conta{ID: 1, GrupoFamiliarID: 1, Nome: "Nubank", Ativo: true}

	t.Run("desativar conta do proprio grupo", func(t *testing.T) {
		// CA-06
		var ativoRecebido *bool
		repo := &mockContasRepo{
			findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return contaExistente, nil },
			setAtivoFn: func(ctx context.Context, id int, ativo bool) error {
				ativoRecebido = &ativo
				return nil
			},
		}
		svc := New(repo)

		conta, err := svc.Desativar(context.Background(), 1, 1)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if conta.Ativo {
			t.Fatal("esperava conta com ativo=false")
		}
		if ativoRecebido == nil || *ativoRecebido != false {
			t.Fatal("esperava chamada ao repositório com ativo=false")
		}
	})

	t.Run("reativar conta do proprio grupo", func(t *testing.T) {
		// CA-07
		repo := &mockContasRepo{
			findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return contaExistente, nil },
			setAtivoFn: func(ctx context.Context, id int, ativo bool) error { return nil },
		}
		svc := New(repo)

		conta, err := svc.Reativar(context.Background(), 1, 1)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !conta.Ativo {
			t.Fatal("esperava conta com ativo=true")
		}
	})

	t.Run("desativar conta de outro grupo familiar nega acesso", func(t *testing.T) {
		// CA-08
		repo := &mockContasRepo{
			findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) { return contaExistente, nil },
		}
		svc := New(repo)

		_, err := svc.Desativar(context.Background(), 2, 1)
		if !errors.Is(err, ErrAcessoNegado) {
			t.Fatalf("esperava ErrAcessoNegado, recebeu %v", err)
		}
	})
}
