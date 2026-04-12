package repository

import (
	"context"
	"fmt"

	"toir-app/internal/models"

	"gorm.io/gorm"
)

// RepairRequestRepository определяет контракт доступа к данным заявок на ремонт.
type RepairRequestRepository interface {
	Create(ctx context.Context, req *models.RepairRequest) error
	FindByID(ctx context.Context, id uint) (*models.RepairRequest, error)
	List(ctx context.Context, page, perPage int, status, priority string, assignedTo *uint) ([]models.RepairRequest, int64, error)
	Update(ctx context.Context, req *models.RepairRequest) error
}

// GormRepairRequestRepository реализует RepairRequestRepository поверх GORM.
type GormRepairRequestRepository struct {
	db *gorm.DB
}

// NewGormRepairRequestRepository создаёт репозиторий заявок на ремонт.
func NewGormRepairRequestRepository(db *gorm.DB) *GormRepairRequestRepository {
	return &GormRepairRequestRepository{db: db}
}

func (r *GormRepairRequestRepository) Create(ctx context.Context, req *models.RepairRequest) error {
	if err := r.db.WithContext(ctx).Create(req).Error; err != nil {
		return fmt.Errorf("failed to create repair request: %w", err)
	}
	return nil
}

func (r *GormRepairRequestRepository) FindByID(ctx context.Context, id uint) (*models.RepairRequest, error) {
	var req models.RepairRequest
	err := r.db.WithContext(ctx).
		Preload("Equipment").
		Preload("Creator").
		Preload("Assignee").
		First(&req, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find repair request by id: %w", err)
	}
	return &req, nil
}

func (r *GormRepairRequestRepository) List(ctx context.Context, page, perPage int, status, priority string, assignedTo *uint) ([]models.RepairRequest, int64, error) {
	var requests []models.RepairRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&models.RepairRequest{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assignedTo != nil {
		query = query.Where("assigned_to = ?", *assignedTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count repair requests: %w", err)
	}

	offset := (page - 1) * perPage
	err := query.
		Preload("Equipment").
		Preload("Creator").
		Preload("Assignee").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&requests).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list repair requests: %w", err)
	}

	return requests, total, nil
}

func (r *GormRepairRequestRepository) Update(ctx context.Context, req *models.RepairRequest) error {
	if err := r.db.WithContext(ctx).Save(req).Error; err != nil {
		return fmt.Errorf("failed to update repair request: %w", err)
	}
	return nil
}
