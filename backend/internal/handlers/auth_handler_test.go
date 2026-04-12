package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"toir-app/internal/handlers"
	"toir-app/internal/models"
	"toir-app/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepo реализует repository.UserRepository для тестов handler.
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *models.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepo) List(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *models.User) error {
	return m.Called(ctx, user).Error(0)
}

func hashTestPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	authSvc := services.NewAuthService(repo, "handler-test-secret")
	h := handlers.NewAuthHandler(authSvc)

	user := &models.User{
		ID:           1,
		Username:     "admin",
		PasswordHash: hashTestPassword(t, "password123"),
		Role:         "admin",
		IsActive:     true,
	}
	repo.On("FindByUsername", mock.Anything, "admin").Return(user, nil)

	e := echo.New()
	body := `{"username":"admin","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "access_token")
	assert.Contains(t, rec.Body.String(), "refresh_token")
	repo.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	authSvc := services.NewAuthService(repo, "handler-test-secret")
	h := handlers.NewAuthHandler(authSvc)

	e := echo.New()
	body := `{invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "error")
}

func TestAuthHandler_Login_EmptyFields(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	authSvc := services.NewAuthService(repo, "handler-test-secret")
	h := handlers.NewAuthHandler(authSvc)

	e := echo.New()
	body := `{"username":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "username and password are required")
}

func TestAuthHandler_Login_WrongCredentials(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	authSvc := services.NewAuthService(repo, "handler-test-secret")
	h := handlers.NewAuthHandler(authSvc)

	repo.On("FindByUsername", mock.Anything, "admin").Return(nil, errors.New("not found"))

	e := echo.New()
	body := `{"username":"admin","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "error")
	repo.AssertExpectations(t)
}
