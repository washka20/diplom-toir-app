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

// Create godoc
// @Summary Создание расписания ТО
// @Description Создаёт новое расписание планового технического обслуживания
// @Tags maintenance-schedules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.MaintenanceSchedule true "Данные расписания"
// @Success 201 {object} response.Response{data=models.MaintenanceSchedule}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /maintenance-schedules [post]
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

// List godoc
// @Summary Список расписаний ТО
// @Description Получение списка расписаний технического обслуживания с пагинацией
// @Tags maintenance-schedules
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param per_page query int false "Элементов на странице" default(20)
// @Success 200 {object} response.Response{data=[]models.MaintenanceSchedule,meta=response.Meta}
// @Failure 500 {object} response.Response
// @Router /maintenance-schedules [get]
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

// Update godoc
// @Summary Обновление расписания ТО
// @Description Обновляет расписание планового технического обслуживания
// @Tags maintenance-schedules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID расписания"
// @Param request body models.MaintenanceSchedule true "Обновлённые данные"
// @Success 200 {object} response.Response{data=models.MaintenanceSchedule}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /maintenance-schedules/{id} [put]
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
