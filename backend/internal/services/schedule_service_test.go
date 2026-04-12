package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"toir-app/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScheduleRepo мок репозитория расписаний.
type MockScheduleRepo struct {
	mock.Mock
}

func (m *MockScheduleRepo) Create(ctx context.Context, s *models.MaintenanceSchedule) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockScheduleRepo) FindByID(ctx context.Context, id uint) (*models.MaintenanceSchedule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MaintenanceSchedule), args.Error(1)
}

func (m *MockScheduleRepo) List(ctx context.Context, page, perPage int) ([]models.MaintenanceSchedule, int64, error) {
	args := m.Called(ctx, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.MaintenanceSchedule), args.Get(1).(int64), args.Error(2)
}

func (m *MockScheduleRepo) Update(ctx context.Context, s *models.MaintenanceSchedule) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockScheduleRepo) FindOverdue(ctx context.Context) ([]models.MaintenanceSchedule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.MaintenanceSchedule), args.Error(1)
}

func TestScheduleService_Create_AutoCalculatesNextDate(t *testing.T) {
	t.Parallel()

	repo := new(MockScheduleRepo)
	svc := NewScheduleService(repo)
	ctx := context.Background()

	schedule := &models.MaintenanceSchedule{
		EquipmentID:  1,
		Type:         "Плановое ТО",
		IntervalDays: 30,
		Description:  "Замена масла",
		CreatedBy:    1,
	}

	repo.On("Create", ctx, mock.AnythingOfType("*models.MaintenanceSchedule")).
		Run(func(args mock.Arguments) {
			s := args.Get(1).(*models.MaintenanceSchedule)
			expected := time.Now().AddDate(0, 0, 30)
			assert.WithinDuration(t, expected, s.NextDate, 2*time.Second)
			assert.True(t, s.IsActive)
		}).
		Return(nil)

	err := svc.Create(ctx, schedule)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestScheduleService_Create_RepoError(t *testing.T) {
	t.Parallel()

	repo := new(MockScheduleRepo)
	svc := NewScheduleService(repo)
	ctx := context.Background()

	schedule := &models.MaintenanceSchedule{
		EquipmentID:  999,
		IntervalDays: 7,
		CreatedBy:    1,
	}

	repo.On("Create", ctx, mock.AnythingOfType("*models.MaintenanceSchedule")).
		Return(fmt.Errorf("failed to create schedule: foreign key violation"))

	err := svc.Create(ctx, schedule)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create schedule")
	repo.AssertExpectations(t)
}

func TestScheduleService_GetByID_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockScheduleRepo)
	svc := NewScheduleService(repo)
	ctx := context.Background()

	expected := &models.MaintenanceSchedule{
		ID:           1,
		EquipmentID:  1,
		Type:         "Плановое ТО",
		IntervalDays: 30,
		IsActive:     true,
	}

	repo.On("FindByID", ctx, uint(1)).Return(expected, nil)

	result, err := svc.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestScheduleService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	repo := new(MockScheduleRepo)
	svc := NewScheduleService(repo)
	ctx := context.Background()

	repo.On("FindByID", ctx, uint(999)).Return(nil, fmt.Errorf("failed to find schedule by id: record not found"))

	result, err := svc.GetByID(ctx, 999)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}

func TestScheduleService_List_WithPagination(t *testing.T) {
	t.Parallel()

	repo := new(MockScheduleRepo)
	svc := NewScheduleService(repo)
	ctx := context.Background()

	items := []models.MaintenanceSchedule{
		{ID: 1, EquipmentID: 1, Type: "Плановое ТО", IntervalDays: 30},
		{ID: 2, EquipmentID: 2, Type: "Диагностика", IntervalDays: 90},
	}

	repo.On("List", ctx, 1, 10).Return(items, int64(2), nil)

	result, total, err := svc.List(ctx, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, items, result)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestScheduleService_Update_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockScheduleRepo)
	svc := NewScheduleService(repo)
	ctx := context.Background()

	schedule := &models.MaintenanceSchedule{
		ID:           1,
		EquipmentID:  1,
		Type:         "Плановое ТО (обновлено)",
		IntervalDays: 14,
		IsActive:     true,
	}

	repo.On("FindByID", ctx, uint(1)).Return(schedule, nil)
	repo.On("Update", ctx, schedule).Return(nil)

	err := svc.Update(ctx, schedule)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestScheduleService_Update_NotFound(t *testing.T) {
	t.Parallel()

	repo := new(MockScheduleRepo)
	svc := NewScheduleService(repo)
	ctx := context.Background()

	schedule := &models.MaintenanceSchedule{
		ID:           999,
		EquipmentID:  1,
		IntervalDays: 14,
	}

	repo.On("FindByID", ctx, uint(999)).Return(nil, fmt.Errorf("failed to find schedule by id: record not found"))

	err := svc.Update(ctx, schedule)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	repo.AssertExpectations(t)
}
