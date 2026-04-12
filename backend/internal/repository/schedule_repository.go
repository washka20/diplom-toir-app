package repository

import (
	"context"
	"fmt"
	"time"

	"toir-app/internal/models"

	"gorm.io/gorm"
)

// ScheduleRepository определяет контракт доступа к данным расписаний ТО.
type ScheduleRepository interface {
	Create(ctx context.Context, s *models.MaintenanceSchedule) error
	FindByID(ctx context.Context, id uint) (*models.MaintenanceSchedule, error)
	List(ctx context.Context, page, perPage int) ([]models.MaintenanceSchedule, int64, error)
	Update(ctx context.Context, s *models.MaintenanceSchedule) error
	FindOverdue(ctx context.Context) ([]models.MaintenanceSchedule, error)
}

// GormScheduleRepository реализует ScheduleRepository поверх GORM.
type GormScheduleRepository struct {
	db *gorm.DB
}

// NewGormScheduleRepository создаёт репозиторий расписаний с подключением к БД.
func NewGormScheduleRepository(db *gorm.DB) *GormScheduleRepository {
	return &GormScheduleRepository{db: db}
}

func (r *GormScheduleRepository) Create(ctx context.Context, s *models.MaintenanceSchedule) error {
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}
	return nil
}

func (r *GormScheduleRepository) FindByID(ctx context.Context, id uint) (*models.MaintenanceSchedule, error) {
	var s models.MaintenanceSchedule
	if err := r.db.WithContext(ctx).Preload("Equipment").Preload("Creator").First(&s, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find schedule by id: %w", err)
	}
	return &s, nil
}

func (r *GormScheduleRepository) List(ctx context.Context, page, perPage int) ([]models.MaintenanceSchedule, int64, error) {
	var items []models.MaintenanceSchedule
	var total int64

	query := r.db.WithContext(ctx).Model(&models.MaintenanceSchedule{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	offset := (page - 1) * perPage
	if err := query.Preload("Equipment").Preload("Creator").Offset(offset).Limit(perPage).Order("id ASC").Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list schedules: %w", err)
	}

	return items, total, nil
}

func (r *GormScheduleRepository) Update(ctx context.Context, s *models.MaintenanceSchedule) error {
	if err := r.db.WithContext(ctx).Save(s).Error; err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}
	return nil
}

func (r *GormScheduleRepository) FindOverdue(ctx context.Context) ([]models.MaintenanceSchedule, error) {
	var items []models.MaintenanceSchedule
	if err := r.db.WithContext(ctx).
		Where("next_date <= ? AND is_active = ?", time.Now(), true).
		Order("next_date ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to find overdue schedules: %w", err)
	}
	return items, nil
}
