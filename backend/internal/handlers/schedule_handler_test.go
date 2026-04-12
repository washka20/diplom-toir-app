package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"toir-app/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScheduleService мок сервиса расписаний для тестирования handler.
type MockScheduleService struct {
	mock.Mock
}

func (m *MockScheduleService) Create(ctx context.Context, s *models.MaintenanceSchedule) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockScheduleService) GetByID(ctx context.Context, id uint) (*models.MaintenanceSchedule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MaintenanceSchedule), args.Error(1)
}

func (m *MockScheduleService) List(ctx context.Context, page, perPage int) ([]models.MaintenanceSchedule, int64, error) {
	args := m.Called(ctx, page, perPage)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.MaintenanceSchedule), args.Get(1).(int64), args.Error(2)
}

func (m *MockScheduleService) Update(ctx context.Context, s *models.MaintenanceSchedule) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func TestScheduleHandler_Create_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"equipment_id":1,"type":"Плановое ТО","interval_days":30,"description":"Замена масла","created_by":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/schedules", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc := new(MockScheduleService)
	svc.On("Create", mock.Anything, mock.AnythingOfType("*models.MaintenanceSchedule")).Return(nil)

	h := NewScheduleHandler(svc)
	err := h.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	svc.AssertExpectations(t)
}

func TestScheduleHandler_Create_ValidationError(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"equipment_id":0,"interval_days":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/schedules", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc := new(MockScheduleService)
	h := NewScheduleHandler(svc)
	err := h.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
	assert.Contains(t, rec.Body.String(), "equipment_id")
}

func TestScheduleHandler_List_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/schedules?page=1&per_page=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	items := []models.MaintenanceSchedule{
		{ID: 1, EquipmentID: 1, Type: "Плановое ТО", IntervalDays: 30},
		{ID: 2, EquipmentID: 2, Type: "Диагностика", IntervalDays: 90},
	}

	svc := new(MockScheduleService)
	svc.On("List", mock.Anything, 1, 10).Return(items, int64(2), nil)

	h := NewScheduleHandler(svc)
	err := h.List(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), `"total":2`)
	svc.AssertExpectations(t)
}

func TestScheduleHandler_Update_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"type":"Плановое ТО (обновлено)","interval_days":14,"is_active":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/schedules/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	svc := new(MockScheduleService)
	svc.On("Update", mock.Anything, mock.AnythingOfType("*models.MaintenanceSchedule")).Return(nil)

	h := NewScheduleHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	svc.AssertExpectations(t)
}

func TestScheduleHandler_Update_InvalidID(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"type":"ТО"}`
	req := httptest.NewRequest(http.MethodPut, "/api/schedules/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	svc := new(MockScheduleService)
	h := NewScheduleHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
}

func TestScheduleHandler_Update_NotFound(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"type":"ТО","interval_days":7}`
	req := httptest.NewRequest(http.MethodPut, "/api/schedules/999", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	svc := new(MockScheduleService)
	svc.On("Update", mock.Anything, mock.AnythingOfType("*models.MaintenanceSchedule")).
		Return(fmt.Errorf("schedule not found"))

	h := NewScheduleHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
	svc.AssertExpectations(t)
}
