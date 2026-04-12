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

// MockWorkOrderService мок сервиса наряд-заказов для тестирования handler.
type MockWorkOrderService struct {
	mock.Mock
}

func (m *MockWorkOrderService) Create(ctx context.Context, wo *models.WorkOrder) error {
	args := m.Called(ctx, wo)
	return args.Error(0)
}

func (m *MockWorkOrderService) GetByID(ctx context.Context, id uint) (*models.WorkOrder, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WorkOrder), args.Error(1)
}

func (m *MockWorkOrderService) List(ctx context.Context, page, perPage int, status string, assignedTo *uint) ([]models.WorkOrder, int64, error) {
	args := m.Called(ctx, page, perPage, status, assignedTo)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.WorkOrder), args.Get(1).(int64), args.Error(2)
}

func (m *MockWorkOrderService) UpdateStatus(ctx context.Context, id uint, newStatus string) (*models.WorkOrder, error) {
	args := m.Called(ctx, id, newStatus)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.WorkOrder), args.Error(1)
}

func TestWorkOrderHandler_Create_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"description":"Ремонт станка","assigned_to":2}`
	req := httptest.NewRequest(http.MethodPost, "/api/work-orders", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc := new(MockWorkOrderService)
	svc.On("Create", mock.Anything, mock.AnythingOfType("*models.WorkOrder")).Return(nil)

	h := NewWorkOrderHandler(svc)
	err := h.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	svc.AssertExpectations(t)
}

func TestWorkOrderHandler_Create_ValidationError(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"assigned_to":0,"description":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/work-orders", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc := new(MockWorkOrderService)
	h := NewWorkOrderHandler(svc)
	err := h.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
	assert.Contains(t, rec.Body.String(), "assigned_to")
}

func TestWorkOrderHandler_List_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/work-orders?page=1&per_page=10&status=pending", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	items := []models.WorkOrder{
		{ID: 1, Description: "Ремонт 1", Status: "pending", AssignedTo: 2},
		{ID: 2, Description: "Ремонт 2", Status: "pending", AssignedTo: 3},
	}

	svc := new(MockWorkOrderService)
	svc.On("List", mock.Anything, 1, 10, "pending", (*uint)(nil)).Return(items, int64(2), nil)

	h := NewWorkOrderHandler(svc)
	err := h.List(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), `"total":2`)
	svc.AssertExpectations(t)
}

func TestWorkOrderHandler_Update_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"status":"in_progress"}`
	req := httptest.NewRequest(http.MethodPut, "/api/work-orders/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	updated := &models.WorkOrder{
		ID:          1,
		Description: "Ремонт станка",
		Status:      "in_progress",
		AssignedTo:  2,
	}

	svc := new(MockWorkOrderService)
	svc.On("UpdateStatus", mock.Anything, uint(1), "in_progress").Return(updated, nil)

	h := NewWorkOrderHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), `"in_progress"`)
	svc.AssertExpectations(t)
}

func TestWorkOrderHandler_Update_InvalidTransition(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"status":"completed"}`
	req := httptest.NewRequest(http.MethodPut, "/api/work-orders/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	svc := new(MockWorkOrderService)
	svc.On("UpdateStatus", mock.Anything, uint(1), "completed").
		Return(nil, fmt.Errorf("invalid status transition from \"pending\" to \"completed\""))

	h := NewWorkOrderHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
	assert.Contains(t, rec.Body.String(), "invalid status transition")
	svc.AssertExpectations(t)
}

func TestWorkOrderHandler_Update_InvalidID(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"status":"in_progress"}`
	req := httptest.NewRequest(http.MethodPut, "/api/work-orders/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	svc := new(MockWorkOrderService)
	h := NewWorkOrderHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
}

func TestWorkOrderHandler_Update_MissingStatus(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{}`
	req := httptest.NewRequest(http.MethodPut, "/api/work-orders/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	svc := new(MockWorkOrderService)
	h := NewWorkOrderHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "status is required")
}
