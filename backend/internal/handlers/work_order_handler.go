package handlers

import (
	"context"
	"net/http"
	"strconv"

	"toir-app/internal/models"
	"toir-app/pkg/response"

	"github.com/labstack/echo/v4"
)

// WorkOrderServiceInterface определяет контракт сервиса наряд-заказов для handler.
type WorkOrderServiceInterface interface {
	Create(ctx context.Context, wo *models.WorkOrder) error
	GetByID(ctx context.Context, id uint) (*models.WorkOrder, error)
	List(ctx context.Context, page, perPage int, status string, assignedTo *uint) ([]models.WorkOrder, int64, error)
	UpdateStatus(ctx context.Context, id uint, newStatus string) (*models.WorkOrder, error)
}

// WorkOrderHandler обрабатывает HTTP-запросы для наряд-заказов.
type WorkOrderHandler struct {
	service WorkOrderServiceInterface
}

// NewWorkOrderHandler создаёт handler наряд-заказов.
func NewWorkOrderHandler(service WorkOrderServiceInterface) *WorkOrderHandler {
	return &WorkOrderHandler{service: service}
}

// UpdateWorkOrderInput представляет тело запроса на обновление статуса наряд-заказа.
type UpdateWorkOrderInput struct {
	Status string `json:"status" example:"in_progress"`
}

// Create godoc
// @Summary Создание наряд-заказа
// @Description Создаёт новый наряд-заказ на выполнение работ
// @Tags work-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.WorkOrder true "Данные наряд-заказа"
// @Success 201 {object} response.Response{data=models.WorkOrder}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /work-orders [post]
func (h *WorkOrderHandler) Create(c echo.Context) error {
	var wo models.WorkOrder
	if err := c.Bind(&wo); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if wo.AssignedTo == 0 || wo.Description == "" {
		return c.JSON(http.StatusBadRequest, response.Error("assigned_to and description are required"))
	}

	if err := h.service.Create(c.Request().Context(), &wo); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(wo))
}

// List godoc
// @Summary Список наряд-заказов
// @Description Получение списка наряд-заказов с пагинацией и фильтрами
// @Tags work-orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param per_page query int false "Элементов на странице" default(20)
// @Param status query string false "Фильтр по статусу"
// @Param assigned_to query int false "Фильтр по исполнителю"
// @Success 200 {object} response.Response{data=[]models.WorkOrder,meta=response.Meta}
// @Failure 500 {object} response.Response
// @Router /work-orders [get]
func (h *WorkOrderHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	status := c.QueryParam("status")

	var assignedTo *uint
	if atStr := c.QueryParam("assigned_to"); atStr != "" {
		if v, err := strconv.ParseUint(atStr, 10, 32); err == nil {
			u := uint(v)
			assignedTo = &u
		}
	}

	items, total, err := h.service.List(c.Request().Context(), page, perPage, status, assignedTo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Paginated(items, page, perPage, total))
}

// Update godoc
// @Summary Обновление наряд-заказа
// @Description Обновляет статус наряд-заказа
// @Tags work-orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID наряд-заказа"
// @Param request body UpdateWorkOrderInput true "Новый статус"
// @Success 200 {object} response.Response{data=models.WorkOrder}
// @Failure 400 {object} response.Response
// @Router /work-orders/{id} [put]
func (h *WorkOrderHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid work order id"))
	}

	var body UpdateWorkOrderInput
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if body.Status == "" {
		return c.JSON(http.StatusBadRequest, response.Error("status is required"))
	}

	wo, err := h.service.UpdateStatus(c.Request().Context(), uint(id), body.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(wo))
}
