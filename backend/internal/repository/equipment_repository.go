package repository

import (
	"context"
	"fmt"

	"toir-app/internal/models"

	"gorm.io/gorm"
)

// EquipmentRepository определяет контракт доступа к данным оборудования.
type EquipmentRepository interface {
	Create(ctx context.Context, equipment *models.Equipment) error
	FindByID(ctx context.Context, id uint) (*models.Equipment, error)
	List(ctx context.Context, page, perPage int, status, location string) ([]models.Equipment, int64, error)
	Update(ctx context.Context, equipment *models.Equipment) error
	Delete(ctx context.Context, id uint) error
}

// GormEquipmentRepository реализует EquipmentRepository поверх GORM.
type GormEquipmentRepository struct {
	db *gorm.DB
}

// NewGormEquipmentRepository создаёт репозиторий оборудования с подключением к БД.
func NewGormEquipmentRepository(db *gorm.DB) *GormEquipmentRepository {
	return &GormEquipmentRepository{db: db}
}

func (r *GormEquipmentRepository) Create(ctx context.Context, equipment *models.Equipment) error {
	if err := r.db.WithContext(ctx).Create(equipment).Error; err != nil {
		return fmt.Errorf("failed to create equipment: %w", err)
	}
	return nil
}

func (r *GormEquipmentRepository) FindByID(ctx context.Context, id uint) (*models.Equipment, error) {
	var equipment models.Equipment
	if err := r.db.WithContext(ctx).First(&equipment, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find equipment by id: %w", err)
	}
	return &equipment, nil
}

func (r *GormEquipmentRepository) List(ctx context.Context, page, perPage int, status, location string) ([]models.Equipment, int64, error) {
	var items []models.Equipment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Equipment{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if location != "" {
		query = query.Where("location = ?", location)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count equipment: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("id ASC").Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list equipment: %w", err)
	}

	return items, total, nil
}

func (r *GormEquipmentRepository) Update(ctx context.Context, equipment *models.Equipment) error {
	if err := r.db.WithContext(ctx).Save(equipment).Error; err != nil {
		return fmt.Errorf("failed to update equipment: %w", err)
	}
	return nil
}

func (r *GormEquipmentRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Equipment{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete equipment: %w", err)
	}
	return nil
}
