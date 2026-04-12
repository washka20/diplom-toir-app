package handlers

import (
	"context"
	"net/http"
	"strconv"

	"toir-app/internal/models"
	"toir-app/pkg/response"

	"github.com/labstack/echo/v4"
)

// ScheduleServiceInterface определяет контракт сервиса расписаний для handler.
type ScheduleServiceInterface interface {
	Create(ctx context.Context, s *models.MaintenanceSchedule) error
	GetByID(ctx context.Context, id uint) (*models.MaintenanceSchedule, error)
	List(ctx context.Context, page, perPage int) ([]models.MaintenanceSchedule, int64, error)
	Update(ctx context.Context, s *models.MaintenanceSchedule) error
}

// ScheduleHandler обрабатывает HTTP-запросы для расписаний ТО.
type ScheduleHandler struct {
	service ScheduleServiceInterface
}

// NewScheduleHandler создаёт handler расписаний.
func NewScheduleHandler(service ScheduleServiceInterface) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

// Create обрабатывает POST /api/schedules.
func (h *ScheduleHandler) Create(c echo.Context) error {
	var s models.MaintenanceSchedule
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if s.EquipmentID == 0 || s.IntervalDays <= 0 {
		return c.JSON(http.StatusBadRequest, response.Error("equipment_id and positive interval_days are required"))
	}

	if err := h.service.Create(c.Request().Context(), &s); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(s))
}

// List обрабатывает GET /api/schedules.
func (h *ScheduleHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	items, total, err := h.service.List(c.Request().Context(), page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Paginated(items, page, perPage, total))
}

// Update обрабатывает PUT /api/schedules/:id.
func (h *ScheduleHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid schedule id"))
	}

	var s models.MaintenanceSchedule
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	s.ID = uint(id)

	if err := h.service.Update(c.Request().Context(), &s); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(s))
}
