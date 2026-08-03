package auth

import (
	"context"
	"errors"
	"expense-control-service/internal/entity"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type mockUsersRepo struct {
	findByEmailFn func(ctx context.Context, email string) (*entity.User, error)
	createFn      func(ctx context.Context, user *entity.User) (*entity.User, error)
}

func (m *mockUsersRepo) ListUserById(context.Context, int) (*entity.User, error) { return nil, nil }
func (m *mockUsersRepo) ListUsers(context.Context) ([]*entity.User, error)       { return nil, nil }
func (m *mockUsersRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.findByEmailFn(ctx, email)
}
func (m *mockUsersRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	return m.createFn(ctx, user)
}

type mockGrupoFamiliarRepo struct {
	createFn func(ctx context.Context) (*entity.GrupoFamiliar, error)
}

func (m *mockGrupoFamiliarRepo) Create(ctx context.Context) (*entity.GrupoFamiliar, error) {
	return m.createFn(ctx)
}

type mockConviteRepo struct {
	createFn          func(ctx context.Context, c *entity.Convite) (*entity.Convite, error)
	findByCodigoFn    func(ctx context.Context, codigo string) (*entity.Convite, error)
	marcarUtilizadoFn func(ctx context.Context, id int, utilizadoEm time.Time) error
}

func (m *mockConviteRepo) Create(ctx context.Context, c *entity.Convite) (*entity.Convite, error) {
	return m.createFn(ctx, c)
}
func (m *mockConviteRepo) FindByCodigo(ctx context.Context, codigo string) (*entity.Convite, error) {
	return m.findByCodigoFn(ctx, codigo)
}
func (m *mockConviteRepo) MarcarUtilizado(ctx context.Context, id int, utilizadoEm time.Time) error {
	return m.marcarUtilizadoFn(ctx, id, utilizadoEm)
}

func TestRegistrar(t *testing.T) {
	codigoValido := "codigo-valido"

	tests := []struct {
		name          string
		nome          string
		email         string
		senha         string
		codigoConvite *string
		usersRepo     *mockUsersRepo
		grupoRepo     *mockGrupoFamiliarRepo
		conviteRepo   *mockConviteRepo
		wantErr       error
		wantGrupoID   int
	}{
		{
			// CA-01
			name:  "sem convite cria grupo familiar novo",
			nome:  "Maristela",
			email: "maristela@example.com",
			senha: "senha1234",
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
				createFn: func(ctx context.Context, u *entity.User) (*entity.User, error) {
					u.ID = 1
					u.CreatedAt = time.Now()
					return u, nil
				},
			},
			grupoRepo: &mockGrupoFamiliarRepo{
				createFn: func(ctx context.Context) (*entity.GrupoFamiliar, error) {
					return &entity.GrupoFamiliar{ID: 10, CriadoEm: time.Now()}, nil
				},
			},
			conviteRepo: &mockConviteRepo{},
			wantGrupoID: 10,
		},
		{
			// CA-02
			name:          "com convite valido vincula ao grupo do convite e marca utilizado",
			nome:          "Trindade",
			email:         "trindade@example.com",
			senha:         "senha1234",
			codigoConvite: &codigoValido,
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
				createFn: func(ctx context.Context, u *entity.User) (*entity.User, error) {
					u.ID = 2
					u.CreatedAt = time.Now()
					return u, nil
				},
			},
			grupoRepo: &mockGrupoFamiliarRepo{
				createFn: func(ctx context.Context) (*entity.GrupoFamiliar, error) {
					t.Fatal("não deveria criar grupo novo quando há convite válido")
					return nil, nil
				},
			},
			conviteRepo: &mockConviteRepo{
				findByCodigoFn: func(ctx context.Context, codigo string) (*entity.Convite, error) {
					return &entity.Convite{ID: 5, Codigo: codigo, GrupoFamiliarID: 10, ExpiraEm: time.Now().Add(24 * time.Hour)}, nil
				},
				marcarUtilizadoFn: func(ctx context.Context, id int, utilizadoEm time.Time) error {
					if id != 5 {
						t.Fatalf("esperava marcar convite id=5, recebeu id=%d", id)
					}
					return nil
				},
			},
			wantGrupoID: 10,
		},
		{
			// CA-03: convite inexistente
			name:          "convite inexistente retorna ErrConvite",
			nome:          "X",
			email:         "x@example.com",
			senha:         "senha1234",
			codigoConvite: &codigoValido,
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
			},
			conviteRepo: &mockConviteRepo{
				findByCodigoFn: func(ctx context.Context, codigo string) (*entity.Convite, error) { return nil, nil },
			},
			wantErr: ErrConvite,
		},
		{
			// CA-03: convite já utilizado
			name:          "convite ja utilizado retorna ErrConvite",
			nome:          "X",
			email:         "x2@example.com",
			senha:         "senha1234",
			codigoConvite: &codigoValido,
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
			},
			conviteRepo: &mockConviteRepo{
				findByCodigoFn: func(ctx context.Context, codigo string) (*entity.Convite, error) {
					usado := time.Now().Add(-time.Hour)
					return &entity.Convite{ID: 1, GrupoFamiliarID: 10, ExpiraEm: time.Now().Add(time.Hour), UtilizadoEm: &usado}, nil
				},
			},
			wantErr: ErrConvite,
		},
		{
			// CA-03: convite expirado
			name:          "convite expirado retorna ErrConvite",
			nome:          "X",
			email:         "x3@example.com",
			senha:         "senha1234",
			codigoConvite: &codigoValido,
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
			},
			conviteRepo: &mockConviteRepo{
				findByCodigoFn: func(ctx context.Context, codigo string) (*entity.Convite, error) {
					return &entity.Convite{ID: 1, GrupoFamiliarID: 10, ExpiraEm: time.Now().Add(-time.Hour)}, nil
				},
			},
			wantErr: ErrConvite,
		},
		{
			// CA-04
			name:  "email ja cadastrado retorna ErrEmailJaCadastrado",
			nome:  "X",
			email: "existente@example.com",
			senha: "senha1234",
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
					return &entity.User{ID: 1, Email: email}, nil
				},
			},
			wantErr: ErrEmailJaCadastrado,
		},
		{
			// CA-05
			name:      "senha curta retorna ErrSenhaCurta",
			nome:      "X",
			email:     "x4@example.com",
			senha:     "1234567",
			usersRepo: &mockUsersRepo{},
			wantErr:   ErrSenhaCurta,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.usersRepo, tt.grupoRepo, tt.conviteRepo, "test-secret")

			user, err := svc.Registrar(context.Background(), tt.nome, tt.email, tt.senha, tt.codigoConvite)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if user.GrupoFamiliarID != tt.wantGrupoID {
				t.Fatalf("esperava grupo_familiar_id=%d, recebeu %d", tt.wantGrupoID, user.GrupoFamiliarID)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	senhaCorreta := "senha1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(senhaCorreta), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("erro ao gerar hash de teste: %v", err)
	}

	usuarioExistente := &entity.User{ID: 1, Email: "user@example.com", SenhaHash: string(hash), GrupoFamiliarID: 10}

	tests := []struct {
		name      string
		email     string
		senha     string
		usersRepo *mockUsersRepo
		wantErr   error
	}{
		{
			// CA-06
			name:  "credenciais corretas retorna token",
			email: usuarioExistente.Email,
			senha: senhaCorreta,
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return usuarioExistente, nil },
			},
		},
		{
			// CA-07
			name:  "senha incorreta retorna ErrCredenciaisInvalidas",
			email: usuarioExistente.Email,
			senha: "senhaErrada",
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return usuarioExistente, nil },
			},
			wantErr: ErrCredenciaisInvalidas,
		},
		{
			name:  "email inexistente retorna ErrCredenciaisInvalidas",
			email: "naoexiste@example.com",
			senha: senhaCorreta,
			usersRepo: &mockUsersRepo{
				findByEmailFn: func(ctx context.Context, email string) (*entity.User, error) { return nil, nil },
			},
			wantErr: ErrCredenciaisInvalidas,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(tt.usersRepo, &mockGrupoFamiliarRepo{}, &mockConviteRepo{}, "test-secret")

			token, expiraEm, err := svc.Login(context.Background(), tt.email, tt.senha)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("esperava erro %v, recebeu %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if token == "" {
				t.Fatal("esperava token não vazio")
			}
			if !expiraEm.After(time.Now()) {
				t.Fatal("esperava expiração no futuro")
			}
			if expiraEm.Sub(time.Now()) > tokenExpiresIn {
				t.Fatalf("expiração maior que %s (RN-12)", tokenExpiresIn)
			}
		})
	}
}

func TestGerarConvite(t *testing.T) {
	conviteRepo := &mockConviteRepo{
		createFn: func(ctx context.Context, c *entity.Convite) (*entity.Convite, error) {
			c.ID = 1
			return c, nil
		},
	}

	svc := New(&mockUsersRepo{}, &mockGrupoFamiliarRepo{}, conviteRepo, "test-secret")

	convite, err := svc.GerarConvite(context.Background(), 10)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if convite.Codigo == "" {
		t.Fatal("esperava código não vazio")
	}
	if convite.GrupoFamiliarID != 10 {
		t.Fatalf("esperava grupo_familiar_id=10, recebeu %d", convite.GrupoFamiliarID)
	}
	if !convite.ExpiraEm.After(time.Now().Add(6*24*time.Hour)) {
		t.Fatal("esperava expiração de aproximadamente 7 dias (RN-09)")
	}
}
