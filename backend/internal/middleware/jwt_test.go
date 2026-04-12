package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"toir-app/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-for-jwt"

func generateTestToken(secret string, userID uint, role string, exp time.Time) string {
	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"role":    role,
		"exp":     jwt.NewNumericDate(exp),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(fmt.Sprintf("failed to sign test token: %v", err))
	}
	return signed
}

func setupEchoContext(token string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	token := generateTestToken(testSecret, 1, "engineer", time.Now().Add(time.Hour))
	c, rec := setupEchoContext(token)

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.NoContent(http.StatusOK)
	}

	mw := middleware.JWTAuth(testSecret)
	err := mw(handler)(c)

	require.NoError(t, err)
	assert.True(t, handlerCalled, "handler should have been called")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, uint(1), c.Get("user_id"))
	assert.Equal(t, "engineer", c.Get("role"))
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	c, rec := setupEchoContext("")

	handler := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}

	mw := middleware.JWTAuth(testSecret)
	err := mw(handler)(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	token := generateTestToken(testSecret, 1, "engineer", time.Now().Add(-time.Hour))
	c, rec := setupEchoContext(token)

	handler := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}

	mw := middleware.JWTAuth(testSecret)
	err := mw(handler)(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_InvalidSignature(t *testing.T) {
	token := generateTestToken("wrong-secret", 1, "engineer", time.Now().Add(time.Hour))
	c, rec := setupEchoContext(token)

	handler := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}

	mw := middleware.JWTAuth(testSecret)
	err := mw(handler)(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
