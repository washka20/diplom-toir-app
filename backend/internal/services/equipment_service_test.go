package services

import (
	"context"
	"fmt"
	"testing"

	"toir-app/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEquipmentRepo мок репозитория оборудования.
type MockEquipmentRepo struct {
	mock.Mock
}

func (m *MockEquipmentRepo) Create(ctx context.Context, equipment *models.Equipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}

func (m *MockEquipmentRepo) FindByID(ctx context.Context, id uint) (*models.Equipment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Equipment), args.Error(1)
}

func (m *MockEquipmentRepo) List(ctx context.Context, page, perPage int, status, location string) ([]models.Equipment, int64, error) {
	args := m.Called(ctx, page, perPage, status, location)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Equipment), args.Get(1).(int64), args.Error(2)
}

func (m *MockEquipmentRepo) Update(ctx context.Context, equipment *models.Equipment) error {
	args := m.Called(ctx, equipment)
	return args.Error(0)
}

func (m *MockEquipmentRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestEquipmentService_Create_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockEquipmentRepo)
	svc := NewEquipmentService(repo)
	ctx := context.Background()

	eq := &models.Equipment{
		Name:            "Станок ЧПУ",
		InventoryNumber: "INV-001",
		Status:          "active",
		Location:        "Цех 1",
	}

	repo.On("Create", ctx, eq).Return(nil)

	err := svc.Create(ctx, eq)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestEquipmentService_Create_DuplicateInventoryNumber(t *testing.T) {
	t.Parallel()

	repo := new(MockEquipmentRepo)
	svc := NewEquipmentService(repo)
	ctx := context.Background()

	eq := &models.Equipment{
		Name:            "Станок ЧПУ",
		InventoryNumber: "INV-001",
	}

	repo.On("Create", ctx, eq).Return(fmt.Errorf("failed to create equipment: duplicate key"))

	err := svc.Create(ctx, eq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create equipment")
	repo.AssertExpectations(t)
}

func TestEquipmentService_GetByID_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockEquipmentRepo)
	svc := NewEquipmentService(repo)
	ctx := context.Background()

	expected := &models.Equipment{
		ID:              1,
		Name:            "Станок ЧПУ",
		InventoryNumber: "INV-001",
		Status:          "active",
	}

	repo.On("FindByID", ctx, uint(1)).Return(expected, nil)

	result, err := svc.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestEquipmentService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	repo := new(MockEquipmentRepo)
	svc := NewEquipmentService(repo)
	ctx := context.Background()

	repo.On("FindByID", ctx, uint(999)).Return(nil, fmt.Errorf("failed to find equipment by id: record not found"))

	result, err := svc.GetByID(ctx, 999)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}

func TestEquipmentService_List_WithPagination(t *testing.T) {
	t.Parallel()

	repo := new(MockEquipmentRepo)
	svc := NewEquipmentService(repo)
	ctx := context.Background()

	tests := []struct {
		name     string
		page     int
		perPage  int
		status   string
		location string
		items    []models.Equipment
		total    int64
	}{
		{
			name:    "первая страница без фильтров",
			page:    1,
			perPage: 10,
			items: []models.Equipment{
				{ID: 1, Name: "Станок 1"},
				{ID: 2, Name: "Станок 2"},
			},
			total: 2,
		},
		{
			name:    "фильтр по статусу",
			page:    1,
			perPage: 10,
			status:  "active",
			items: []models.Equipment{
				{ID: 1, Name: "Станок 1", Status: "active"},
			},
			total: 1,
		},
		{
			name:     "фильтр по локации",
			page:     1,
			perPage:  10,
			location: "Цех 1",
			items: []models.Equipment{
				{ID: 1, Name: "Станок 1", Location: "Цех 1"},
			},
			total: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			localRepo := new(MockEquipmentRepo)
			localSvc := NewEquipmentService(localRepo)

			localRepo.On("List", ctx, tc.page, tc.perPage, tc.status, tc.location).
				Return(tc.items, tc.total, nil)

			items, total, err := localSvc.List(ctx, tc.page, tc.perPage, tc.status, tc.location)

			assert.NoError(t, err)
			assert.Equal(t, tc.items, items)
			assert.Equal(t, tc.total, total)
			localRepo.AssertExpectations(t)
		})
	}

	_ = repo
	_ = svc
}

func TestEquipmentService_Update_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockEquipmentRepo)
	svc := NewEquipmentService(repo)
	ctx := context.Background()

	eq := &models.Equipment{
		ID:              1,
		Name:            "Станок ЧПУ (обновлён)",
		InventoryNumber: "INV-001",
		Status:          "maintenance",
	}

	repo.On("FindByID", ctx, uint(1)).Return(eq, nil)
	repo.On("Update", ctx, eq).Return(nil)

	err := svc.Update(ctx, eq)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestEquipmentService_Delete_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockEquipmentRepo)
	svc := NewEquipmentService(repo)
	ctx := context.Background()

	repo.On("FindByID", ctx, uint(1)).Return(&models.Equipment{ID: 1}, nil)
	repo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.Delete(ctx, 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
