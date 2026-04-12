package services_test

import (
	"context"
	"errors"
	"testing"

	"toir-app/internal/models"
	"toir-app/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepo реализует repository.UserRepository через testify/mock.
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

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func TestAuthService_Login_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	svc := services.NewAuthService(repo, "test-secret")

	user := &models.User{
		ID:           1,
		Username:     "admin",
		PasswordHash: hashPassword(t, "correctpassword"),
		Role:         "admin",
		IsActive:     true,
	}
	repo.On("FindByUsername", mock.Anything, "admin").Return(user, nil)

	tokens, err := svc.Login(context.Background(), "admin", "correctpassword")

	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	svc := services.NewAuthService(repo, "test-secret")

	user := &models.User{
		ID:           1,
		Username:     "admin",
		PasswordHash: hashPassword(t, "correctpassword"),
		Role:         "admin",
		IsActive:     true,
	}
	repo.On("FindByUsername", mock.Anything, "admin").Return(user, nil)

	tokens, err := svc.Login(context.Background(), "admin", "wrongpassword")

	require.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "invalid credentials")
	repo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	svc := services.NewAuthService(repo, "test-secret")

	repo.On("FindByUsername", mock.Anything, "unknown").Return(nil, errors.New("failed to find user by username: record not found"))

	tokens, err := svc.Login(context.Background(), "unknown", "anypassword")

	require.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "invalid credentials")
	repo.AssertExpectations(t)
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	t.Parallel()
	repo := new(MockUserRepo)
	svc := services.NewAuthService(repo, "test-secret")

	user := &models.User{
		ID:           2,
		Username:     "inactive",
		PasswordHash: hashPassword(t, "password"),
		Role:         "engineer",
		IsActive:     false,
	}
	repo.On("FindByUsername", mock.Anything, "inactive").Return(user, nil)

	tokens, err := svc.Login(context.Background(), "inactive", "password")

	require.Error(t, err)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "user account is deactivated")
	repo.AssertExpectations(t)
}
