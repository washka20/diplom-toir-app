package repository

import (
	"context"
	"fmt"

	"toir-app/internal/models"

	"gorm.io/gorm"
)

// UserRepository определяет контракт доступа к данным пользователей.
type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
}

// GormUserRepository реализует UserRepository поверх GORM.
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository создаёт репозиторий пользователей с подключением к БД.
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}
	return &user, nil
}

func (r *GormUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *GormUserRepository) List(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
