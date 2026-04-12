package services

import (
	"context"
	"fmt"
	"time"

	"toir-app/internal/models"
	"toir-app/internal/repository"
)

// allowedTransitions определяет допустимые переходы статусов заявки.
var allowedTransitions = map[string][]string{
	"new":           {"assigned"},
	"assigned":      {"in_progress"},
	"in_progress":   {"waiting_parts", "completed"},
	"waiting_parts": {"in_progress"},
	"completed":     {"closed"},
}

// RepairRequestService инкапсулирует бизнес-логику заявок на ремонт.
type RepairRequestService struct {
	repo repository.RepairRequestRepository
}

// NewRepairRequestService создаёт сервис заявок на ремонт.
func NewRepairRequestService(repo repository.RepairRequestRepository) *RepairRequestService {
	return &RepairRequestService{repo: repo}
}

// Create создаёт новую заявку, устанавливая статус "new" и привязывая к пользователю.
func (s *RepairRequestService) Create(ctx context.Context, req *models.RepairRequest, userID uint) error {
	req.Status = "new"
	req.CreatedBy = userID

	if err := s.repo.Create(ctx, req); err != nil {
		return fmt.Errorf("failed to create repair request: %w", err)
	}
	return nil
}

// GetByID возвращает заявку по ID с предзагрузкой связей.
func (s *RepairRequestService) GetByID(ctx context.Context, id uint) (*models.RepairRequest, error) {
	req, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repair request not found: %w", err)
	}
	return req, nil
}

// List возвращает список заявок с пагинацией и фильтрацией.
func (s *RepairRequestService) List(ctx context.Context, page, perPage int, status, priority string, assignedTo *uint) ([]models.RepairRequest, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	return s.repo.List(ctx, page, perPage, status, priority, assignedTo)
}

// Assign назначает техника на заявку и переводит статус в "assigned".
func (s *RepairRequestService) Assign(ctx context.Context, id uint, technicianID uint) error {
	req, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repair request not found: %w", err)
	}

	if !isTransitionAllowed(req.Status, "assigned") {
		return fmt.Errorf("invalid status transition from %q to %q", req.Status, "assigned")
	}

	req.AssignedTo = &technicianID
	req.Status = "assigned"

	if err := s.repo.Update(ctx, req); err != nil {
		return fmt.Errorf("failed to assign repair request: %w", err)
	}
	return nil
}

// UpdateStatus обновляет статус заявки с валидацией допустимого перехода.
func (s *RepairRequestService) UpdateStatus(ctx context.Context, id uint, newStatus string) error {
	req, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repair request not found: %w", err)
	}

	if !isTransitionAllowed(req.Status, newStatus) {
		return fmt.Errorf("invalid status transition from %q to %q", req.Status, newStatus)
	}

	req.Status = newStatus

	if newStatus == "completed" {
		now := time.Now()
		req.CompletedAt = &now
	}

	if err := s.repo.Update(ctx, req); err != nil {
		return fmt.Errorf("failed to update repair request status: %w", err)
	}
	return nil
}

func isTransitionAllowed(from, to string) bool {
	targets, ok := allowedTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}
