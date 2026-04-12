package repository

import (
	"context"
	"fmt"

	"toir-app/internal/models"

	"gorm.io/gorm"
)

// WorkOrderRepository определяет контракт доступа к данным наряд-заказов.
type WorkOrderRepository interface {
	Create(ctx context.Context, wo *models.WorkOrder) error
	FindByID(ctx context.Context, id uint) (*models.WorkOrder, error)
	List(ctx context.Context, page, perPage int, status string, assignedTo *uint) ([]models.WorkOrder, int64, error)
	Update(ctx context.Context, wo *models.WorkOrder) error
}

// GormWorkOrderRepository реализует WorkOrderRepository поверх GORM.
type GormWorkOrderRepository struct {
	db *gorm.DB
}

// NewGormWorkOrderRepository создаёт репозиторий наряд-заказов с подключением к БД.
func NewGormWorkOrderRepository(db *gorm.DB) *GormWorkOrderRepository {
	return &GormWorkOrderRepository{db: db}
}

func (r *GormWorkOrderRepository) Create(ctx context.Context, wo *models.WorkOrder) error {
	if err := r.db.WithContext(ctx).Create(wo).Error; err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}
	return nil
}

func (r *GormWorkOrderRepository) FindByID(ctx context.Context, id uint) (*models.WorkOrder, error) {
	var wo models.WorkOrder
	if err := r.db.WithContext(ctx).First(&wo, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find work order by id: %w", err)
	}
	return &wo, nil
}

func (r *GormWorkOrderRepository) List(ctx context.Context, page, perPage int, status string, assignedTo *uint) ([]models.WorkOrder, int64, error) {
	var items []models.WorkOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&models.WorkOrder{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if assignedTo != nil {
		query = query.Where("assigned_to = ?", *assignedTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count work orders: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Order("id ASC").Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list work orders: %w", err)
	}

	return items, total, nil
}

func (r *GormWorkOrderRepository) Update(ctx context.Context, wo *models.WorkOrder) error {
	if err := r.db.WithContext(ctx).Save(wo).Error; err != nil {
		return fmt.Errorf("failed to update work order: %w", err)
	}
	return nil
}
