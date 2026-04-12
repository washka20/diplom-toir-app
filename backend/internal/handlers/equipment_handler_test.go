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

// MockEquipmentService мок сервиса оборудования для тестирования handler.
type MockEquipmentService struct {
	mock.Mock
}

func (m *MockEquipmentService) Create(ctx context.Context, eq *models.Equipment) error {
	args := m.Called(ctx, eq)
	return args.Error(0)
}

func (m *MockEquipmentService) GetByID(ctx context.Context, id uint) (*models.Equipment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Equipment), args.Error(1)
}

func (m *MockEquipmentService) List(ctx context.Context, page, perPage int, status, location string) ([]models.Equipment, int64, error) {
	args := m.Called(ctx, page, perPage, status, location)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Equipment), args.Get(1).(int64), args.Error(2)
}

func (m *MockEquipmentService) Update(ctx context.Context, eq *models.Equipment) error {
	args := m.Called(ctx, eq)
	return args.Error(0)
}

func (m *MockEquipmentService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestEquipmentHandler_Create_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"name":"Станок ЧПУ","inventory_number":"INV-001","status":"active","location":"Цех 1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/equipment", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	svc := new(MockEquipmentService)
	svc.On("Create", mock.Anything, mock.AnythingOfType("*models.Equipment")).Return(nil)

	h := NewEquipmentHandler(svc)
	err := h.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	svc.AssertExpectations(t)
}

func TestEquipmentHandler_List_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/equipment?page=1&per_page=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	items := []models.Equipment{
		{ID: 1, Name: "Станок 1", Status: "active"},
		{ID: 2, Name: "Станок 2", Status: "active"},
	}

	svc := new(MockEquipmentService)
	svc.On("List", mock.Anything, 1, 10, "", "").Return(items, int64(2), nil)

	h := NewEquipmentHandler(svc)
	err := h.List(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), `"total":2`)
	svc.AssertExpectations(t)
}

func TestEquipmentHandler_GetByID_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/equipment/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	expected := &models.Equipment{
		ID:              1,
		Name:            "Станок ЧПУ",
		InventoryNumber: "INV-001",
		Status:          "active",
	}

	svc := new(MockEquipmentService)
	svc.On("GetByID", mock.Anything, uint(1)).Return(expected, nil)

	h := NewEquipmentHandler(svc)
	err := h.GetByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), "Станок ЧПУ")
	svc.AssertExpectations(t)
}

func TestEquipmentHandler_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/equipment/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	svc := new(MockEquipmentService)
	svc.On("GetByID", mock.Anything, uint(999)).Return(nil, fmt.Errorf("equipment not found"))

	h := NewEquipmentHandler(svc)
	err := h.GetByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
	svc.AssertExpectations(t)
}

func TestEquipmentHandler_Update_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	body := `{"name":"Станок ЧПУ (обновлён)","inventory_number":"INV-001","status":"maintenance"}`
	req := httptest.NewRequest(http.MethodPut, "/api/equipment/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	svc := new(MockEquipmentService)
	svc.On("Update", mock.Anything, mock.AnythingOfType("*models.Equipment")).Return(nil)

	h := NewEquipmentHandler(svc)
	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	svc.AssertExpectations(t)
}

func TestEquipmentHandler_Delete_Success(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/equipment/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	svc := new(MockEquipmentService)
	svc.On("Delete", mock.Anything, uint(1)).Return(nil)

	h := NewEquipmentHandler(svc)
	err := h.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	svc.AssertExpectations(t)
}
