package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"toir-app/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDashboardService мок сервиса дашборда для тестирования handler.
type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetMetrics(ctx context.Context) (*models.DashboardMetrics, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardMetrics), args.Error(1)
}

func TestDashboardHandler_GetMetrics_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expected := &models.DashboardMetrics{
		OpenRequests:       5,
		OverdueSchedules:   2,
		AvgRepairTimeHrs:   3.5,
		CompletedThisMonth: 12,
		TotalEquipment:     50,
		ActiveEquipment:    42,
	}

	svc := new(MockDashboardService)
	svc.On("GetMetrics", mock.Anything).Return(expected, nil)

	h := NewDashboardHandler(svc)
	err := h.GetMetrics(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), `"open_requests":5`)
	assert.Contains(t, rec.Body.String(), `"overdue_schedules":2`)
	assert.Contains(t, rec.Body.String(), `"avg_repair_time_hrs":3.5`)
	assert.Contains(t, rec.Body.String(), `"completed_this_month":12`)
	assert.Contains(t, rec.Body.String(), `"total_equipment":50`)
	assert.Contains(t, rec.Body.String(), `"active_equipment":42`)
	svc.AssertExpectations(t)
}

func TestDashboardHandler_GetMetrics_ServiceError(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/metrics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc := new(MockDashboardService)
	svc.On("GetMetrics", mock.Anything).Return(nil, fmt.Errorf("database connection failed"))

	h := NewDashboardHandler(svc)
	err := h.GetMetrics(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
	svc.AssertExpectations(t)
}
