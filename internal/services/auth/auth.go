package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"expense-control-service/internal/entity"
	conviteRepo "expense-control-service/internal/repositories/convite"
	grupoFamiliarRepo "expense-control-service/internal/repositories/grupo_familiar"
	usersRepo "expense-control-service/internal/repositories/users"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenExpiresIn  = time.Hour        // RN-12
	conviteValidity = 7 * 24 * time.Hour // RN-09
	senhaMinLength  = 8                // RN-05
)

type Auth interface {
	Registrar(ctx context.Context, nome, email, senha string, codigoConvite *string) (*entity.User, error)
	Login(ctx context.Context, email, senha string) (token string, expiraEm time.Time, err error)
	GerarConvite(ctx context.Context, grupoFamiliarID int) (*entity.Convite, error)
}

type auth struct {
	usersRepo         usersRepo.Users
	grupoFamiliarRepo grupoFamiliarRepo.GrupoFamiliar
	conviteRepo       conviteRepo.Convite
	jwtSecret         []byte
}

func New(
	usersRepo usersRepo.Users,
	grupoFamiliarRepo grupoFamiliarRepo.GrupoFamiliar,
	conviteRepo conviteRepo.Convite,
	jwtSecret string,
) Auth {
	return &auth{
		usersRepo:         usersRepo,
		grupoFamiliarRepo: grupoFamiliarRepo,
		conviteRepo:       conviteRepo,
		jwtSecret:         []byte(jwtSecret),
	}
}

func (a *auth) Registrar(ctx context.Context, nome, email, senha string, codigoConvite *string) (*entity.User, error) {
	if len(senha) < senhaMinLength {
		return nil, ErrSenhaCurta
	}

	existente, err := a.usersRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, ErrEmailJaCadastrado
	}

	grupoFamiliarID, err := a.resolverGrupoFamiliar(ctx, codigoConvite)
	if err != nil {
		return nil, err
	}

	senhaHash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	user := &entity.User{
		Nome:            nome,
		Email:           email,
		SenhaHash:       string(senhaHash),
		GrupoFamiliarID: grupoFamiliarID,
	}

	created, err := a.usersRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, usersRepo.ErrEmailDuplicado) {
			// RN-04: outra requisição concorrente cadastrou o mesmo e-mail entre o
			// FindByEmail acima e este Create — a constraint UNIQUE do banco pegou a
			// corrida (TOCTOU), aqui só traduzimos para o erro de domínio esperado.
			return nil, ErrEmailJaCadastrado
		}
		return nil, err
	}

	return created, nil
}

// resolverGrupoFamiliar aplica RN-01/RN-02: sem código de convite, cria um grupo
// familiar novo; com código válido, vincula ao grupo do convite e o marca como
// utilizado (RN-08), sem criar grupo novo.
func (a *auth) resolverGrupoFamiliar(ctx context.Context, codigoConvite *string) (int, error) {
	if codigoConvite == nil || *codigoConvite == "" {
		grupo, err := a.grupoFamiliarRepo.Create(ctx)
		if err != nil {
			return 0, err
		}
		return grupo.ID, nil
	}

	c, err := a.conviteRepo.FindByCodigo(ctx, *codigoConvite)
	if err != nil {
		return 0, err
	}
	if c == nil {
		return 0, fmt.Errorf("%w: código de convite não encontrado", ErrConvite)
	}
	if c.Utilizado() {
		return 0, fmt.Errorf("%w: código de convite já utilizado", ErrConvite)
	}
	if c.Expirado(time.Now()) {
		return 0, fmt.Errorf("%w: código de convite expirado", ErrConvite)
	}

	if err := a.conviteRepo.MarcarUtilizado(ctx, c.ID, time.Now()); err != nil {
		if errors.Is(err, conviteRepo.ErrJaUtilizado) {
			// RN-08: outra requisição concorrente consumiu o convite entre o FindByCodigo
			// acima e este UPDATE — fecha a janela de corrida (TOCTOU).
			return 0, fmt.Errorf("%w: código de convite já utilizado", ErrConvite)
		}
		return 0, err
	}

	return c.GrupoFamiliarID, nil
}

func (a *auth) Login(ctx context.Context, email, senha string) (string, time.Time, error) {
	user, err := a.usersRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, err
	}
	if user == nil {
		return "", time.Time{}, ErrCredenciaisInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.SenhaHash), []byte(senha)); err != nil {
		return "", time.Time{}, ErrCredenciaisInvalidas
	}

	return a.generateToken(user)
}

// generateToken cria um JWT com os claims usuario_id e grupo_familiar_id — RN-11/RN-12.
func (a *auth) generateToken(user *entity.User) (string, time.Time, error) {
	expiraEm := time.Now().Add(tokenExpiresIn)

	claims := jwt.MapClaims{
		"usuario_id":        user.ID,
		"grupo_familiar_id": user.GrupoFamiliarID,
		"exp":               expiraEm.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("erro ao assinar token: %w", err)
	}

	return signed, expiraEm, nil
}

func (a *auth) GerarConvite(ctx context.Context, grupoFamiliarID int) (*entity.Convite, error) {
	codigo, err := gerarCodigo()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar código de convite: %w", err)
	}

	convite := &entity.Convite{
		Codigo:          codigo,
		GrupoFamiliarID: grupoFamiliarID,
		ExpiraEm:        time.Now().Add(conviteValidity),
	}

	return a.conviteRepo.Create(ctx, convite)
}

func gerarCodigo() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
