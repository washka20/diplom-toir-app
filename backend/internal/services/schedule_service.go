package services

import (
	"context"
	"fmt"
	"time"

	"toir-app/internal/models"
	"toir-app/internal/repository"
)

// ScheduleService инкапсулирует бизнес-логику работы с расписаниями ТО.
type ScheduleService struct {
	repo repository.ScheduleRepository
}

// NewScheduleService создаёт сервис расписаний.
func NewScheduleService(repo repository.ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

// Create создаёт новое расписание ТО, автоматически рассчитывая next_date.
func (s *ScheduleService) Create(ctx context.Context, schedule *models.MaintenanceSchedule) error {
	schedule.NextDate = time.Now().AddDate(0, 0, schedule.IntervalDays)
	schedule.IsActive = true

	if err := s.repo.Create(ctx, schedule); err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}
	return nil
}

// GetByID возвращает расписание по ID.
func (s *ScheduleService) GetByID(ctx context.Context, id uint) (*models.MaintenanceSchedule, error) {
	schedule, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("schedule not found: %w", err)
	}
	return schedule, nil
}

// List возвращает список расписаний с пагинацией.
func (s *ScheduleService) List(ctx context.Context, page, perPage int) ([]models.MaintenanceSchedule, int64, error) {
	items, total, err := s.repo.List(ctx, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list schedules: %w", err)
	}
	return items, total, nil
}

// Update обновляет расписание ТО.
func (s *ScheduleService) Update(ctx context.Context, schedule *models.MaintenanceSchedule) error {
	if _, err := s.repo.FindByID(ctx, schedule.ID); err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}
	if err := s.repo.Update(ctx, schedule); err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}
	return nil
}
