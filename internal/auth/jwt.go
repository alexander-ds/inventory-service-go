package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Claims refleja el contrato actual del JWT de auth-service.
// Los roles no se incluyen: son una decisión futura (docs/authentication.md).
type Claims struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	jwt.RegisteredClaims
}

// Validator verifica tokens HS256 firmados por auth-service.
type Validator struct {
	secret []byte
}

func NewValidator(secret string) *Validator {
	return &Validator{secret: []byte(secret)}
}

// Parse valida algoritmo, firma y expiración. Devuelve los claims verificados.
func (v *Validator) Parse(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token vacío")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("algoritmo no permitido: %s", token.Method.Alg())
			}
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("algoritmo no permitido: %s", token.Method.Alg())
			}
			return v.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
