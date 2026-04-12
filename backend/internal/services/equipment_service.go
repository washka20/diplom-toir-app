package services

import (
	"context"
	"fmt"

	"toir-app/internal/models"
	"toir-app/internal/repository"
)

// EquipmentService инкапсулирует бизнес-логику работы с оборудованием.
type EquipmentService struct {
	repo repository.EquipmentRepository
}

// NewEquipmentService создаёт сервис оборудования.
func NewEquipmentService(repo repository.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: repo}
}

// Create создаёт новое оборудование.
func (s *EquipmentService) Create(ctx context.Context, eq *models.Equipment) error {
	if err := s.repo.Create(ctx, eq); err != nil {
		return fmt.Errorf("failed to create equipment: %w", err)
	}
	return nil
}

// GetByID возвращает оборудование по ID.
func (s *EquipmentService) GetByID(ctx context.Context, id uint) (*models.Equipment, error) {
	eq, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("equipment not found: %w", err)
	}
	return eq, nil
}

// List возвращает список оборудования с пагинацией и фильтрами.
func (s *EquipmentService) List(ctx context.Context, page, perPage int, status, location string) ([]models.Equipment, int64, error) {
	items, total, err := s.repo.List(ctx, page, perPage, status, location)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list equipment: %w", err)
	}
	return items, total, nil
}

// Update обновляет данные оборудования.
func (s *EquipmentService) Update(ctx context.Context, eq *models.Equipment) error {
	if _, err := s.repo.FindByID(ctx, eq.ID); err != nil {
		return fmt.Errorf("equipment not found: %w", err)
	}
	if err := s.repo.Update(ctx, eq); err != nil {
		return fmt.Errorf("failed to update equipment: %w", err)
	}
	return nil
}

// Delete удаляет оборудование по ID.
func (s *EquipmentService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return fmt.Errorf("equipment not found: %w", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete equipment: %w", err)
	}
	return nil
}
