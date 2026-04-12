package services

import (
	"context"
	"fmt"

	"toir-app/internal/models"
	"toir-app/internal/repository"
)

// woAllowedTransitions определяет допустимые переходы статусов наряд-заказа.
var woAllowedTransitions = map[string]string{
	"pending":     "in_progress",
	"in_progress": "completed",
}

// WorkOrderService инкапсулирует бизнес-логику работы с наряд-заказами.
type WorkOrderService struct {
	repo repository.WorkOrderRepository
}

// NewWorkOrderService создаёт сервис наряд-заказов.
func NewWorkOrderService(repo repository.WorkOrderRepository) *WorkOrderService {
	return &WorkOrderService{repo: repo}
}

// Create создаёт новый наряд-заказ со статусом pending.
func (s *WorkOrderService) Create(ctx context.Context, wo *models.WorkOrder) error {
	wo.Status = "pending"
	if err := s.repo.Create(ctx, wo); err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}
	return nil
}

// GetByID возвращает наряд-заказ по ID.
func (s *WorkOrderService) GetByID(ctx context.Context, id uint) (*models.WorkOrder, error) {
	wo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("work order not found: %w", err)
	}
	return wo, nil
}

// List возвращает список наряд-заказов с пагинацией и фильтрами.
func (s *WorkOrderService) List(ctx context.Context, page, perPage int, status string, assignedTo *uint) ([]models.WorkOrder, int64, error) {
	items, total, err := s.repo.List(ctx, page, perPage, status, assignedTo)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list work orders: %w", err)
	}
	return items, total, nil
}

// UpdateStatus обновляет статус наряд-заказа с валидацией перехода.
func (s *WorkOrderService) UpdateStatus(ctx context.Context, id uint, newStatus string) (*models.WorkOrder, error) {
	wo, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("work order not found: %w", err)
	}

	allowed, ok := woAllowedTransitions[wo.Status]
	if !ok || allowed != newStatus {
		return nil, fmt.Errorf("invalid status transition from %q to %q", wo.Status, newStatus)
	}

	wo.Status = newStatus
	if err := s.repo.Update(ctx, wo); err != nil {
		return nil, fmt.Errorf("failed to update work order status: %w", err)
	}
	return wo, nil
}
