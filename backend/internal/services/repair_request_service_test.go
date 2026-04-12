package services_test

import (
	"context"
	"errors"
	"testing"

	"toir-app/internal/models"
	"toir-app/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepairRequestRepo реализует repository.RepairRequestRepository через testify/mock.
type MockRepairRequestRepo struct {
	mock.Mock
}

func (m *MockRepairRequestRepo) Create(ctx context.Context, req *models.RepairRequest) error {
	return m.Called(ctx, req).Error(0)
}

func (m *MockRepairRequestRepo) FindByID(ctx context.Context, id uint) (*models.RepairRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RepairRequest), args.Error(1)
}

func (m *MockRepairRequestRepo) List(ctx context.Context, page, perPage int, status, priority string, assignedTo *uint) ([]models.RepairRequest, int64, error) {
	args := m.Called(ctx, page, perPage, status, priority, assignedTo)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.RepairRequest), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepairRequestRepo) Update(ctx context.Context, req *models.RepairRequest) error {
	return m.Called(ctx, req).Error(0)
}

func TestRepairRequestService_Create_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)

	repo.On("Create", mock.Anything, mock.MatchedBy(func(req *models.RepairRequest) bool {
		return req.Status == "new" && req.CreatedBy == uint(1)
	})).Return(nil)

	req := &models.RepairRequest{
		EquipmentID: 10,
		Title:       "Сломался станок",
		Description: "Не запускается",
		Priority:    "high",
	}

	err := svc.Create(context.Background(), req, 1)

	require.NoError(t, err)
	assert.Equal(t, "new", req.Status)
	assert.Equal(t, uint(1), req.CreatedBy)
	repo.AssertExpectations(t)
}

func TestRepairRequestService_GetByID_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)

	expected := &models.RepairRequest{
		ID:          1,
		Title:       "Заявка",
		Status:      "new",
		EquipmentID: 5,
		CreatedBy:   2,
	}
	repo.On("FindByID", mock.Anything, uint(1)).Return(expected, nil)

	result, err := svc.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestRepairRequestService_GetByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)

	repo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("failed to find repair request by id: record not found"))

	result, err := svc.GetByID(context.Background(), 999)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}

func TestRepairRequestService_List_WithFilters(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)

	assignedTo := uint(3)
	expected := []models.RepairRequest{
		{ID: 1, Status: "assigned", Priority: "high", AssignedTo: &assignedTo},
		{ID: 2, Status: "assigned", Priority: "high", AssignedTo: &assignedTo},
	}
	repo.On("List", mock.Anything, 1, 20, "assigned", "high", &assignedTo).Return(expected, int64(2), nil)

	results, total, err := svc.List(context.Background(), 1, 20, "assigned", "high", &assignedTo)

	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestRepairRequestService_Assign_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)

	existing := &models.RepairRequest{
		ID:     1,
		Status: "new",
		Title:  "Заявка",
	}
	repo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(req *models.RepairRequest) bool {
		return req.Status == "assigned" && req.AssignedTo != nil && *req.AssignedTo == uint(5)
	})).Return(nil)

	err := svc.Assign(context.Background(), 1, 5)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRepairRequestService_UpdateStatus_ValidTransition(t *testing.T) {
	t.Parallel()

	transitions := []struct {
		name      string
		from      string
		to        string
	}{
		{"new to assigned", "new", "assigned"},
		{"assigned to in_progress", "assigned", "in_progress"},
		{"in_progress to completed", "in_progress", "completed"},
		{"in_progress to waiting_parts", "in_progress", "waiting_parts"},
		{"waiting_parts to in_progress", "waiting_parts", "in_progress"},
		{"completed to closed", "completed", "closed"},
	}

	for _, tc := range transitions {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := new(MockRepairRequestRepo)
			svc := services.NewRepairRequestService(repo)

			existing := &models.RepairRequest{
				ID:     1,
				Status: tc.from,
			}
			repo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
			repo.On("Update", mock.Anything, mock.MatchedBy(func(req *models.RepairRequest) bool {
				return req.Status == tc.to
			})).Return(nil)

			err := svc.UpdateStatus(context.Background(), 1, tc.to)

			require.NoError(t, err)
			repo.AssertExpectations(t)
		})
	}
}

func TestRepairRequestService_UpdateStatus_InvalidTransition(t *testing.T) {
	t.Parallel()

	invalidTransitions := []struct {
		name string
		from string
		to   string
	}{
		{"new to completed", "new", "completed"},
		{"new to in_progress", "new", "in_progress"},
		{"assigned to completed", "assigned", "completed"},
		{"completed to new", "completed", "new"},
		{"closed to new", "closed", "new"},
		{"new to closed", "new", "closed"},
	}

	for _, tc := range invalidTransitions {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := new(MockRepairRequestRepo)
			svc := services.NewRepairRequestService(repo)

			existing := &models.RepairRequest{
				ID:     1,
				Status: tc.from,
			}
			repo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)

			err := svc.UpdateStatus(context.Background(), 1, tc.to)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid status transition")
			repo.AssertExpectations(t)
		})
	}
}
