package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"toir-app/internal/handlers"
	"toir-app/internal/models"
	"toir-app/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepairRequestRepo реализует repository.RepairRequestRepository для тестов handler.
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

func TestRepairRequestHandler_Create_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)
	h := handlers.NewRepairRequestHandler(svc)

	repo.On("Create", mock.Anything, mock.MatchedBy(func(req *models.RepairRequest) bool {
		return req.Status == "new" && req.CreatedBy == uint(1) && req.Title == "Сломался станок"
	})).Return(nil)

	e := echo.New()
	body := `{"equipment_id":10,"title":"Сломался станок","description":"Не запускается","priority":"high"}`
	req := httptest.NewRequest(http.MethodPost, "/api/repair-requests", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", uint(1))

	err := h.Create(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	repo.AssertExpectations(t)
}

func TestRepairRequestHandler_List_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)
	h := handlers.NewRepairRequestHandler(svc)

	expected := []models.RepairRequest{
		{ID: 1, Title: "Заявка 1", Status: "new", Priority: "high"},
		{ID: 2, Title: "Заявка 2", Status: "new", Priority: "low"},
	}
	repo.On("List", mock.Anything, 1, 20, "new", "", (*uint)(nil)).Return(expected, int64(2), nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/repair-requests?page=1&per_page=20&status=new", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.List(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), `"total":2`)
	repo.AssertExpectations(t)
}

func TestRepairRequestHandler_GetByID_Success(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)
	h := handlers.NewRepairRequestHandler(svc)

	expected := &models.RepairRequest{
		ID:          1,
		Title:       "Заявка",
		Status:      "new",
		Priority:    "high",
		EquipmentID: 5,
		CreatedBy:   2,
	}
	repo.On("FindByID", mock.Anything, uint(1)).Return(expected, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/repair-requests/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.GetByID(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	assert.Contains(t, rec.Body.String(), `"Заявка"`)
	repo.AssertExpectations(t)
}

func TestRepairRequestHandler_GetByID_NotFound(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)
	h := handlers.NewRepairRequestHandler(svc)

	repo.On("FindByID", mock.Anything, uint(999)).Return(nil, errors.New("not found"))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/repair-requests/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err := h.GetByID(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":false`)
	repo.AssertExpectations(t)
}

func TestRepairRequestHandler_UpdateStatus(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)
	h := handlers.NewRepairRequestHandler(svc)

	existing := &models.RepairRequest{
		ID:     1,
		Status: "assigned",
		Title:  "Заявка",
	}
	repo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(req *models.RepairRequest) bool {
		return req.Status == "in_progress"
	})).Return(nil)

	e := echo.New()
	body := `{"status":"in_progress"}`
	req := httptest.NewRequest(http.MethodPut, "/api/repair-requests/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.Update(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	repo.AssertExpectations(t)
}

func TestRepairRequestHandler_Update_Assign(t *testing.T) {
	t.Parallel()
	repo := new(MockRepairRequestRepo)
	svc := services.NewRepairRequestService(repo)
	h := handlers.NewRepairRequestHandler(svc)

	existing := &models.RepairRequest{
		ID:     1,
		Status: "new",
		Title:  "Заявка",
	}
	repo.On("FindByID", mock.Anything, uint(1)).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(req *models.RepairRequest) bool {
		return req.Status == "assigned" && req.AssignedTo != nil && *req.AssignedTo == uint(5)
	})).Return(nil)

	e := echo.New()
	body := `{"assigned_to":5}`
	req := httptest.NewRequest(http.MethodPut, "/api/repair-requests/1", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.Update(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"success":true`)
	repo.AssertExpectations(t)
}
