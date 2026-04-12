package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDashboardQuerier мок для DashboardDBQuerier.
type MockDashboardQuerier struct {
	mock.Mock
}

func (m *MockDashboardQuerier) CountOpenRequests(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardQuerier) CountOverdueSchedules(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardQuerier) AvgRepairTimeThisMonth(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDashboardQuerier) CountCompletedThisMonth(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardQuerier) CountTotalEquipment(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDashboardQuerier) CountActiveEquipment(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestDashboardService_GetMetrics_Success(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(5), nil)
	querier.On("CountOverdueSchedules", ctx).Return(int64(2), nil)
	querier.On("AvgRepairTimeThisMonth", ctx).Return(3.456, nil)
	querier.On("CountCompletedThisMonth", ctx).Return(int64(12), nil)
	querier.On("CountTotalEquipment", ctx).Return(int64(50), nil)
	querier.On("CountActiveEquipment", ctx).Return(int64(42), nil)

	metrics, err := svc.GetMetrics(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), metrics.OpenRequests)
	assert.Equal(t, int64(2), metrics.OverdueSchedules)
	assert.Equal(t, 3.46, metrics.AvgRepairTimeHrs)
	assert.Equal(t, int64(12), metrics.CompletedThisMonth)
	assert.Equal(t, int64(50), metrics.TotalEquipment)
	assert.Equal(t, int64(42), metrics.ActiveEquipment)
	querier.AssertExpectations(t)
}

func TestDashboardService_GetMetrics_ZeroAvgRepairTime(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(0), nil)
	querier.On("CountOverdueSchedules", ctx).Return(int64(0), nil)
	querier.On("AvgRepairTimeThisMonth", ctx).Return(0.0, nil)
	querier.On("CountCompletedThisMonth", ctx).Return(int64(0), nil)
	querier.On("CountTotalEquipment", ctx).Return(int64(0), nil)
	querier.On("CountActiveEquipment", ctx).Return(int64(0), nil)

	metrics, err := svc.GetMetrics(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), metrics.OpenRequests)
	assert.Equal(t, 0.0, metrics.AvgRepairTimeHrs)
	querier.AssertExpectations(t)
}

func TestDashboardService_GetMetrics_OpenRequestsError(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(0), fmt.Errorf("db error"))

	metrics, err := svc.GetMetrics(ctx)

	assert.Nil(t, metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count open requests")
	querier.AssertExpectations(t)
}

func TestDashboardService_GetMetrics_OverdueSchedulesError(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(5), nil)
	querier.On("CountOverdueSchedules", ctx).Return(int64(0), fmt.Errorf("db error"))

	metrics, err := svc.GetMetrics(ctx)

	assert.Nil(t, metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count overdue schedules")
	querier.AssertExpectations(t)
}

func TestDashboardService_GetMetrics_AvgRepairTimeError(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(5), nil)
	querier.On("CountOverdueSchedules", ctx).Return(int64(2), nil)
	querier.On("AvgRepairTimeThisMonth", ctx).Return(0.0, fmt.Errorf("db error"))

	metrics, err := svc.GetMetrics(ctx)

	assert.Nil(t, metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get avg repair time")
	querier.AssertExpectations(t)
}

func TestDashboardService_GetMetrics_CompletedError(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(5), nil)
	querier.On("CountOverdueSchedules", ctx).Return(int64(2), nil)
	querier.On("AvgRepairTimeThisMonth", ctx).Return(3.5, nil)
	querier.On("CountCompletedThisMonth", ctx).Return(int64(0), fmt.Errorf("db error"))

	metrics, err := svc.GetMetrics(ctx)

	assert.Nil(t, metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count completed this month")
	querier.AssertExpectations(t)
}

func TestDashboardService_GetMetrics_TotalEquipmentError(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(5), nil)
	querier.On("CountOverdueSchedules", ctx).Return(int64(2), nil)
	querier.On("AvgRepairTimeThisMonth", ctx).Return(3.5, nil)
	querier.On("CountCompletedThisMonth", ctx).Return(int64(12), nil)
	querier.On("CountTotalEquipment", ctx).Return(int64(0), fmt.Errorf("db error"))

	metrics, err := svc.GetMetrics(ctx)

	assert.Nil(t, metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count total equipment")
	querier.AssertExpectations(t)
}

func TestDashboardService_GetMetrics_ActiveEquipmentError(t *testing.T) {
	t.Parallel()

	querier := new(MockDashboardQuerier)
	svc := NewDashboardService(querier)
	ctx := context.Background()

	querier.On("CountOpenRequests", ctx).Return(int64(5), nil)
	querier.On("CountOverdueSchedules", ctx).Return(int64(2), nil)
	querier.On("AvgRepairTimeThisMonth", ctx).Return(3.5, nil)
	querier.On("CountCompletedThisMonth", ctx).Return(int64(12), nil)
	querier.On("CountTotalEquipment", ctx).Return(int64(50), nil)
	querier.On("CountActiveEquipment", ctx).Return(int64(0), fmt.Errorf("db error"))

	metrics, err := svc.GetMetrics(ctx)

	assert.Nil(t, metrics)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to count active equipment")
	querier.AssertExpectations(t)
}
