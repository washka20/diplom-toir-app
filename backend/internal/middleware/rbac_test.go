package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"toir-app/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireRole_Allowed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("role", "engineer")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.NoContent(http.StatusOK)
	}

	mw := middleware.RequireRole("engineer", "admin")
	err := mw(handler)(c)

	require.NoError(t, err)
	assert.True(t, handlerCalled, "handler should have been called")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireRole_Forbidden(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("role", "operator")

	handler := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}

	mw := middleware.RequireRole("engineer", "admin")
	err := mw(handler)(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestRequireRole_MissingRole(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}

	mw := middleware.RequireRole("engineer", "admin")
	err := mw(handler)(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
