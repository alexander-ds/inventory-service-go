package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "unit_test_secret"

func signToken(t *testing.T, secret string, claims jwt.MapClaims, method jwt.SigningMethod) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("no se pudo firmar el token de prueba: %v", err)
	}
	return signed
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub":   "3f1b0c2e-0000-4000-8000-000000000001",
		"email": "admin@inventory.local",
		"name":  "Admin",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
}

func TestParseValidToken(t *testing.T) {
	token := signToken(t, testSecret, validClaims(), jwt.SigningMethodHS256)

	claims, err := NewValidator(testSecret).Parse(token)
	if err != nil {
		t.Fatalf("esperaba token válido, obtuve: %v", err)
	}

	if claims.Subject != "3f1b0c2e-0000-4000-8000-000000000001" {
		t.Errorf("sub incorrecto: %s", claims.Subject)
	}
	if claims.Email != "admin@inventory.local" || claims.Name != "Admin" {
		t.Errorf("claims de identidad incorrectos: %+v", claims)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	token := signToken(t, testSecret, validClaims(), jwt.SigningMethodHS256)

	if _, err := NewValidator("otro_secret").Parse(token); err == nil {
		t.Fatal("esperaba error con un secret distinto")
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	claims := validClaims()
	claims["exp"] = time.Now().Add(-time.Hour).Unix()
	token := signToken(t, testSecret, claims, jwt.SigningMethodHS256)

	if _, err := NewValidator(testSecret).Parse(token); err == nil {
		t.Fatal("esperaba error con token expirado")
	}
}

func TestParseRejectsOtherAlgorithm(t *testing.T) {
	token := signToken(t, testSecret, validClaims(), jwt.SigningMethodHS512)

	if _, err := NewValidator(testSecret).Parse(token); err == nil {
		t.Fatal("esperaba error con algoritmo distinto de HS256")
	}
}

func TestParseRejectsEmptyToken(t *testing.T) {
	if _, err := NewValidator(testSecret).Parse(""); err == nil {
		t.Fatal("esperaba error con token vacío")
	}
}
