package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextKeyUsuarioID       = "usuario_id"
	ContextKeyGrupoFamiliarID = "grupo_familiar_id"
)

// AuthMiddleware valida o JWT no header Authorization: Bearer <token> — RN-13.
// Disponibiliza usuario_id e grupo_familiar_id no contexto Gin para uso downstream (RN-14).
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	secret := []byte(jwtSecret)

	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "token ausente ou mal formatado"})
			return
		}

		claims, err := parseToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": nil, "error": mapTokenErrorMessage(err)})
			return
		}

		if !setContextValues(c, claims) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"data": nil, "error": "token com claims inválidas"})
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

func parseToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.MapClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return *claims, nil
}

func mapTokenErrorMessage(err error) string {
	if errors.Is(err, jwt.ErrTokenExpired) {
		return "token expirado"
	}
	return "token inválido"
}

// setContextValues extrai usuario_id e grupo_familiar_id dos claims. Retorna false
// se algum dos dois campos obrigatórios (RN-11) estiver ausente ou em formato inesperado.
func setContextValues(c *gin.Context, claims jwt.MapClaims) bool {
	usuarioID, ok := claims["usuario_id"].(float64)
	if !ok {
		return false
	}

	grupoFamiliarID, ok := claims["grupo_familiar_id"].(float64)
	if !ok {
		return false
	}

	c.Set(ContextKeyUsuarioID, int(usuarioID))
	c.Set(ContextKeyGrupoFamiliarID, int(grupoFamiliarID))
	return true
}
