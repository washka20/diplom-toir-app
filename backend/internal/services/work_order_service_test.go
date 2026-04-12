package services

import (
	"context"
	"fmt"
	"testing"

	"toir-app/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorkOrderRepo мок репозитория наряд-заказов.
type MockWorkOrderRepo struct {
	mock.Mock
}

func (m *MockWorkOrderRepo) Create(ctx context.Context, wo *models.WorkOrder) error {
	args := m.Called(ctx, wo)
	return args.Error(0)
}

func (m *MockWorkOrderRepo) FindByID(ctx context.Context, id uint) (*models.WorkOrder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WorkOrder), args.Error(1)
}

func (m *MockWorkOrderRepo) List(ctx context.Context, page, perPage int, status string, assignedTo *uint) ([]models.WorkOrder, int64, error) {
	args := m.Called(ctx, page, perPage, status, assignedTo)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.WorkOrder), args.Get(1).(int64), args.Error(2)
}

func (m *MockWorkOrderRepo) Update(ctx context.Context, wo *models.WorkOrder) error {
	args := m.Called(ctx, wo)
	return args.Error(0)
}

func TestWorkOrderService_Create_SetsPendingStatus(t *testing.T) {
	t.Parallel()

	repo := new(MockWorkOrderRepo)
	svc := NewWorkOrderService(repo)
	ctx := context.Background()

	wo := &models.WorkOrder{
		Description: "Ремонт станка",
		AssignedTo:  2,
	}

	repo.On("Create", ctx, mock.AnythingOfType("*models.WorkOrder")).
		Run(func(args mock.Arguments) {
			order := args.Get(1).(*models.WorkOrder)
			assert.Equal(t, "pending", order.Status)
		}).
		Return(nil)

	err := svc.Create(ctx, wo)

	assert.NoError(t, err)
	assert.Equal(t, "pending", wo.Status)
	repo.AssertExpectations(t)
}

func TestWorkOrderService_Create_RepoError(t *testing.T) {
	t.Parallel()

	repo := new(MockWorkOrderRepo)
	svc := NewWorkOrderService(repo)
	ctx := context.Background()

	wo := &models.WorkOrder{
		Description: "Ремонт",
		AssignedTo:  999,
	}

	repo.On("Create", ctx, mock.AnythingOfType("*models.WorkOrder")).
		Return(fmt.Errorf("failed to create work order: foreign key violation"))

	err := svc.Create(ctx, wo)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create work order")
	repo.AssertExpectations(t)
}

func TestWorkOrderService_UpdateStatus_PendingToInProgress(t *testing.T) {
	t.Parallel()

	repo := new(MockWorkOrderRepo)
	svc := NewWorkOrderService(repo)
	ctx := context.Background()

	existing := &models.WorkOrder{
		ID:          1,
		Description: "Ремонт станка",
		Status:      "pending",
		AssignedTo:  2,
	}

	repo.On("FindByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*models.WorkOrder")).Return(nil)

	result, err := svc.UpdateStatus(ctx, 1, "in_progress")

	assert.NoError(t, err)
	assert.Equal(t, "in_progress", result.Status)
	repo.AssertExpectations(t)
}

func TestWorkOrderService_UpdateStatus_InProgressToCompleted(t *testing.T) {
	t.Parallel()

	repo := new(MockWorkOrderRepo)
	svc := NewWorkOrderService(repo)
	ctx := context.Background()

	existing := &models.WorkOrder{
		ID:          1,
		Description: "Ремонт станка",
		Status:      "in_progress",
		AssignedTo:  2,
	}

	repo.On("FindByID", ctx, uint(1)).Return(existing, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*models.WorkOrder")).Return(nil)

	result, err := svc.UpdateStatus(ctx, 1, "completed")

	assert.NoError(t, err)
	assert.Equal(t, "completed", result.Status)
	repo.AssertExpectations(t)
}

func TestWorkOrderService_UpdateStatus_InvalidTransition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		currentStatus string
		newStatus     string
	}{
		{
			name:          "pending напрямую в completed",
			currentStatus: "pending",
			newStatus:     "completed",
		},
		{
			name:          "completed назад в pending",
			currentStatus: "completed",
			newStatus:     "pending",
		},
		{
			name:          "in_progress назад в pending",
			currentStatus: "in_progress",
			newStatus:     "pending",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := new(MockWorkOrderRepo)
			svc := NewWorkOrderService(repo)
			ctx := context.Background()

			existing := &models.WorkOrder{
				ID:     1,
				Status: tc.currentStatus,
			}

			repo.On("FindByID", ctx, uint(1)).Return(existing, nil)

			result, err := svc.UpdateStatus(ctx, 1, tc.newStatus)

			assert.Nil(t, result)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid status transition")
			repo.AssertExpectations(t)
		})
	}
}

func TestWorkOrderService_UpdateStatus_NotFound(t *testing.T) {
	t.Parallel()

	repo := new(MockWorkOrderRepo)
	svc := NewWorkOrderService(repo)
	ctx := context.Background()

	repo.On("FindByID", ctx, uint(999)).Return(nil, fmt.Errorf("failed to find work order by id: record not found"))

	result, err := svc.UpdateStatus(ctx, 999, "in_progress")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}

func TestWorkOrderService_List_WithFilters(t *testing.T) {
	t.Parallel()

	repo := new(MockWorkOrderRepo)
	svc := NewWorkOrderService(repo)
	ctx := context.Background()

	assignedTo := uint(2)
	items := []models.WorkOrder{
		{ID: 1, Description: "Ремонт 1", Status: "pending", AssignedTo: 2},
	}

	repo.On("List", ctx, 1, 10, "pending", &assignedTo).Return(items, int64(1), nil)

	result, total, err := svc.List(ctx, 1, 10, "pending", &assignedTo)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestWorkOrderService_GetByID_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockWorkOrderRepo)
	svc := NewWorkOrderService(repo)
	ctx := context.Background()

	expected := &models.WorkOrder{
		ID:          1,
		Description: "Ремонт станка",
		Status:      "pending",
		AssignedTo:  2,
	}

	repo.On("FindByID", ctx, uint(1)).Return(expected, nil)

	result, err := svc.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}
