package launches

import (
	"context"
	"errors"
	"testing"
	"time"

	"expense-control-service/internal/entity"
)

type mockLaunchesRepo struct {
	listFn       func(ctx context.Context, grupoFamiliarID int, mes, ano *int) ([]*entity.Launches, error)
	findByIDFn   func(ctx context.Context, id int) (*entity.Launches, error)
	createFn     func(ctx context.Context, l *entity.Launches) (*entity.Launches, error)
	updateFn     func(ctx context.Context, l *entity.Launches) (*entity.Launches, error)
	softDeleteFn func(ctx context.Context, id int) error
}

func (m *mockLaunchesRepo) ListLaunches(ctx context.Context, grupoFamiliarID int, mes, ano *int) ([]*entity.Launches, error) {
	return m.listFn(ctx, grupoFamiliarID, mes, ano)
}
func (m *mockLaunchesRepo) FindByID(ctx context.Context, id int) (*entity.Launches, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockLaunchesRepo) Create(ctx context.Context, l *entity.Launches) (*entity.Launches, error) {
	return m.createFn(ctx, l)
}
func (m *mockLaunchesRepo) Update(ctx context.Context, l *entity.Launches) (*entity.Launches, error) {
	return m.updateFn(ctx, l)
}
func (m *mockLaunchesRepo) SoftDelete(ctx context.Context, id int) error {
	return m.softDeleteFn(ctx, id)
}

type mockContasRepo struct {
	findByIDFn func(ctx context.Context, id int) (*entity.Conta, error)
}

func (m *mockContasRepo) Create(ctx context.Context, c *entity.Conta) (*entity.Conta, error) {
	return nil, nil
}
func (m *mockContasRepo) FindByID(ctx context.Context, id int) (*entity.Conta, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockContasRepo) ListByGrupoFamiliar(ctx context.Context, grupoFamiliarID int) ([]*entity.Conta, error) {
	return nil, nil
}
func (m *mockContasRepo) Update(ctx context.Context, c *entity.Conta) (*entity.Conta, error) {
	return nil, nil
}
func (m *mockContasRepo) SetAtivo(ctx context.Context, id int, ativo bool) error { return nil }

type mockLaunchCategoryRepo struct {
	findByIDFn func(ctx context.Context, id int) (*entity.LaunchCategory, error)
}

func (m *mockLaunchCategoryRepo) ListLaunchCategorys(ctx context.Context) ([]*entity.LaunchCategory, error) {
	return nil, nil
}
func (m *mockLaunchCategoryRepo) FindByID(ctx context.Context, id int) (*entity.LaunchCategory, error) {
	return m.findByIDFn(ctx, id)
}

type mockPaymentMethodsRepo struct {
	existsFn func(ctx context.Context, id int) (bool, error)
}

func (m *mockPaymentMethodsRepo) ListPaymentMethodById(ctx context.Context, id int) (*entity.PaymentMethods, error) {
	return nil, nil
}
func (m *mockPaymentMethodsRepo) Exists(ctx context.Context, id int) (bool, error) {
	return m.existsFn(ctx, id)
}

const (
	grupoFamiliar1       = 1
	grupoFamiliar2       = 2
	contaAtivaID         = 1
	contaInativaID       = 2
	contaOutroGrupo      = 3
	categoriaDespesaID   = 1
	categoriaReceitaID   = 2
	formaPagamentoValida = 1
)

func repos() (*mockContasRepo, *mockLaunchCategoryRepo, *mockPaymentMethodsRepo) {
	contas := &mockContasRepo{
		findByIDFn: func(ctx context.Context, id int) (*entity.Conta, error) {
			switch id {
			case contaAtivaID:
				return &entity.Conta{ID: contaAtivaID, GrupoFamiliarID: grupoFamiliar1, Ativo: true}, nil
			case contaInativaID:
				return &entity.Conta{ID: contaInativaID, GrupoFamiliarID: grupoFamiliar1, Ativo: false}, nil
			case contaOutroGrupo:
				return &entity.Conta{ID: contaOutroGrupo, GrupoFamiliarID: grupoFamiliar2, Ativo: true}, nil
			default:
				return nil, nil
			}
		},
	}
	categorias := &mockLaunchCategoryRepo{
		findByIDFn: func(ctx context.Context, id int) (*entity.LaunchCategory, error) {
			switch id {
			case categoriaDespesaID:
				return &entity.LaunchCategory{ID: categoriaDespesaID, TipoLancamentoID: uint(entity.LauchTypeSaida.ID)}, nil
			case categoriaReceitaID:
				return &entity.LaunchCategory{ID: categoriaReceitaID, TipoLancamentoID: uint(entity.LauchTypeEntrada.ID)}, nil
			default:
				return nil, nil
			}
		},
	}
	formasPagamento := &mockPaymentMethodsRepo{
		existsFn: func(ctx context.Context, id int) (bool, error) {
			return id == formaPagamentoValida, nil
		},
	}
	return contas, categorias, formasPagamento
}

func TestCriar(t *testing.T) {
	tests := []struct {
		name             string
		grupoFamiliarID  int
		contaID          int
		tipo             string
		categoriaID      int
		formaPagamentoID int
		valor            float64
		wantErr          error
	}{
		{
			// CA-01
			name: "lancamento valido criado com sucesso", grupoFamiliarID: grupoFamiliar1,
			contaID: contaAtivaID, tipo: TipoDespesa, categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: 50,
		},
		{
			// CA-02: conta inexistente
			name: "conta inexistente retorna ErrContaInvalida", grupoFamiliarID: grupoFamiliar1,
			contaID: 999, tipo: TipoDespesa, categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: 50, wantErr: ErrContaInvalida,
		},
		{
			// CA-02: conta de outro grupo
			name: "conta de outro grupo retorna ErrContaInvalida", grupoFamiliarID: grupoFamiliar1,
			contaID: contaOutroGrupo, tipo: TipoDespesa, categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: 50, wantErr: ErrContaInvalida,
		},
		{
			// CA-03: conta desativada
			name: "conta desativada retorna ErrContaInvalida", grupoFamiliarID: grupoFamiliar1,
			contaID: contaInativaID, tipo: TipoDespesa, categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: 50, wantErr: ErrContaInvalida,
		},
		{
			// CA-04
			name: "tipo invalido retorna ErrTipoInvalido", grupoFamiliarID: grupoFamiliar1,
			contaID: contaAtivaID, tipo: "invalido", categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: 50, wantErr: ErrTipoInvalido,
		},
		{
			// CA-05: categoria incompatível com o tipo
			name: "categoria de despesa em lancamento receita retorna ErrCategoriaInvalida", grupoFamiliarID: grupoFamiliar1,
			contaID: contaAtivaID, tipo: TipoReceita, categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: 50, wantErr: ErrCategoriaInvalida,
		},
		{
			name: "categoria inexistente retorna ErrCategoriaInvalida", grupoFamiliarID: grupoFamiliar1,
			contaID: contaAtivaID, tipo: TipoDespesa, categoriaID: 999,
			formaPagamentoID: formaPagamentoValida, valor: 50, wantErr: ErrCategoriaInvalida,
		},
		{
			name: "forma de pagamento invalida retorna ErrFormaPagamentoInvalida", grupoFamiliarID: grupoFamiliar1,
			contaID: contaAtivaID, tipo: TipoDespesa, categoriaID: categoriaDespesaID,
			formaPagamentoID: 999, valor: 50, wantErr: ErrFormaPagamentoInvalida,
		},
		{
			// CA-06
			name: "valor zero retorna ErrValorInvalido", grupoFamiliarID: grupoFamiliar1,
			contaID: contaAtivaID, tipo: TipoDespesa, categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: 0, wantErr: ErrValorInvalido,
		},
		{
			name: "valor negativo retorna ErrValorInvalido", grupoFamiliarID: grupoFamiliar1,
			contaID: contaAtivaID, tipo: TipoDespesa, categoriaID: categoriaDespesaID,
			formaPagamentoID: formaPagamentoValida, valor: -10, wantErr: ErrValorInvalido,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contas, categorias, formasPagamento := repos()
			launchesRepo := &mockLaunchesRepo{
				createFn: func(ctx context.Context, l *entity.Launches) (*entity.Launches, error) {
					l.ID = 1
					return l, nil
				},
			}
			svc := New(launchesRepo, contas, categorias, formasPagamento)

			launch, err := svc.Criar(context.Background(), tt.grupoFamiliarID, 10, tt.contaID, tt.tipo, tt.categoriaID, tt.formaPagamentoID, "desc", tt.valor, time.Now())

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if launch.GrupoFamiliarID != tt.grupoFamiliarID {
				t.Fatalf("esperava grupo_familiar_id=%d, recebeu %d", tt.grupoFamiliarID, launch.GrupoFamiliarID)
			}
			if launch.Usuario != 10 {
				t.Fatalf("esperava usuario=10 (RN-02), recebeu %d", launch.Usuario)
			}
		})
	}
}

func TestListLaunches(t *testing.T) {
	// CA-07/CA-08
	contas, categorias, formasPagamento := repos()
	launchesRepo := &mockLaunchesRepo{
		listFn: func(ctx context.Context, grupoFamiliarID int, mes, ano *int) ([]*entity.Launches, error) {
			if grupoFamiliarID != grupoFamiliar1 {
				t.Fatalf("esperava grupo_familiar_id=%d, recebeu %d", grupoFamiliar1, grupoFamiliarID)
			}
			return []*entity.Launches{{ID: 1, GrupoFamiliarID: grupoFamiliar1}}, nil
		},
	}
	svc := New(launchesRepo, contas, categorias, formasPagamento)

	result, err := svc.ListLaunches(context.Background(), grupoFamiliar1, nil, nil)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("esperava 1 lançamento, recebeu %d", len(result))
	}
}

func TestAtualizar(t *testing.T) {
	lancamentoExistente := &entity.Launches{ID: 1, GrupoFamiliarID: grupoFamiliar1, Usuario: 10}

	t.Run("edita lancamento do proprio grupo com sucesso", func(t *testing.T) {
		// CA-09
		contas, categorias, formasPagamento := repos()
		launchesRepo := &mockLaunchesRepo{
			findByIDFn: func(ctx context.Context, id int) (*entity.Launches, error) { return lancamentoExistente, nil },
			updateFn:   func(ctx context.Context, l *entity.Launches) (*entity.Launches, error) { return l, nil },
		}
		svc := New(launchesRepo, contas, categorias, formasPagamento)

		_, err := svc.Atualizar(context.Background(), grupoFamiliar1, 1, contaAtivaID, TipoDespesa, categoriaDespesaID, formaPagamentoValida, "nova desc", 100, time.Now())
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
	})

	t.Run("lancamento inexistente retorna ErrLancamentoNaoEncontrado", func(t *testing.T) {
		contas, categorias, formasPagamento := repos()
		launchesRepo := &mockLaunchesRepo{
			findByIDFn: func(ctx context.Context, id int) (*entity.Launches, error) { return nil, nil },
		}
		svc := New(launchesRepo, contas, categorias, formasPagamento)

		_, err := svc.Atualizar(context.Background(), grupoFamiliar1, 999, contaAtivaID, TipoDespesa, categoriaDespesaID, formaPagamentoValida, "d", 100, time.Now())
		if !errors.Is(err, ErrLancamentoNaoEncontrado) {
			t.Fatalf("esperava ErrLancamentoNaoEncontrado, recebeu %v", err)
		}
	})

	t.Run("lancamento de outro grupo familiar retorna ErrAcessoNegado", func(t *testing.T) {
		// CA-11
		contas, categorias, formasPagamento := repos()
		launchesRepo := &mockLaunchesRepo{
			findByIDFn: func(ctx context.Context, id int) (*entity.Launches, error) { return lancamentoExistente, nil },
		}
		svc := New(launchesRepo, contas, categorias, formasPagamento)

		_, err := svc.Atualizar(context.Background(), grupoFamiliar2, 1, contaAtivaID, TipoDespesa, categoriaDespesaID, formaPagamentoValida, "d", 100, time.Now())
		if !errors.Is(err, ErrAcessoNegado) {
			t.Fatalf("esperava ErrAcessoNegado, recebeu %v", err)
		}
	})
}

func TestExcluir(t *testing.T) {
	lancamentoExistente := &entity.Launches{ID: 1, GrupoFamiliarID: grupoFamiliar1}

	t.Run("exclui lancamento do proprio grupo (soft-delete)", func(t *testing.T) {
		// CA-10
		contas, categorias, formasPagamento := repos()
		var softDeleteChamado bool
		launchesRepo := &mockLaunchesRepo{
			findByIDFn:   func(ctx context.Context, id int) (*entity.Launches, error) { return lancamentoExistente, nil },
			softDeleteFn: func(ctx context.Context, id int) error { softDeleteChamado = true; return nil },
		}
		svc := New(launchesRepo, contas, categorias, formasPagamento)

		if err := svc.Excluir(context.Background(), grupoFamiliar1, 1); err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if !softDeleteChamado {
			t.Fatal("esperava chamada ao soft-delete do repositório")
		}
	})

	t.Run("exclui lancamento de outro grupo familiar nega acesso", func(t *testing.T) {
		// CA-11
		contas, categorias, formasPagamento := repos()
		launchesRepo := &mockLaunchesRepo{
			findByIDFn: func(ctx context.Context, id int) (*entity.Launches, error) { return lancamentoExistente, nil },
		}
		svc := New(launchesRepo, contas, categorias, formasPagamento)

		err := svc.Excluir(context.Background(), grupoFamiliar2, 1)
		if !errors.Is(err, ErrAcessoNegado) {
			t.Fatalf("esperava ErrAcessoNegado, recebeu %v", err)
		}
	})
}
