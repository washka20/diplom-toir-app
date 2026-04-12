package handlers

import (
	"net/http"
	"strconv"

	"toir-app/internal/services"
	"toir-app/pkg/response"

	"toir-app/internal/models"

	"github.com/labstack/echo/v4"
)

// RepairRequestHandler обрабатывает HTTP-запросы заявок на ремонт.
type RepairRequestHandler struct {
	service *services.RepairRequestService
}

// NewRepairRequestHandler создаёт handler заявок на ремонт.
func NewRepairRequestHandler(service *services.RepairRequestService) *RepairRequestHandler {
	return &RepairRequestHandler{service: service}
}

// CreateRepairRequestInput представляет тело запроса на создание заявки.
type CreateRepairRequestInput struct {
	EquipmentID uint   `json:"equipment_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// UpdateRepairRequestInput представляет тело запроса на обновление заявки.
type UpdateRepairRequestInput struct {
	Status     string `json:"status"`
	AssignedTo *uint  `json:"assigned_to"`
}

// Create обрабатывает POST /api/repair-requests.
func (h *RepairRequestHandler) Create(c echo.Context) error {
	var input CreateRepairRequestInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if input.Title == "" || input.Priority == "" || input.EquipmentID == 0 {
		return c.JSON(http.StatusBadRequest, response.Error("title, priority and equipment_id are required"))
	}

	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, response.Error("unauthorized"))
	}

	req := &models.RepairRequest{
		EquipmentID: input.EquipmentID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
	}

	if err := h.service.Create(c.Request().Context(), req, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(req))
}

// List обрабатывает GET /api/repair-requests.
func (h *RepairRequestHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	status := c.QueryParam("status")
	priority := c.QueryParam("priority")

	var assignedTo *uint
	if v := c.QueryParam("assigned_to"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			uid := uint(id)
			assignedTo = &uid
		}
	}

	results, total, err := h.service.List(c.Request().Context(), page, perPage, status, priority, assignedTo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Paginated(results, page, perPage, total))
}

// GetByID обрабатывает GET /api/repair-requests/:id.
func (h *RepairRequestHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid id"))
	}

	req, err := h.service.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("repair request not found"))
	}

	return c.JSON(http.StatusOK, response.Success(req))
}

// Update обрабатывает PUT /api/repair-requests/:id.
func (h *RepairRequestHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid id"))
	}

	var input UpdateRepairRequestInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	ctx := c.Request().Context()
	uid := uint(id)

	if input.AssignedTo != nil {
		if err := h.service.Assign(ctx, uid, *input.AssignedTo); err != nil {
			return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		}
		return c.JSON(http.StatusOK, response.Success(map[string]string{"message": "repair request assigned"}))
	}

	if input.Status != "" {
		if err := h.service.UpdateStatus(ctx, uid, input.Status); err != nil {
			return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		}
		return c.JSON(http.StatusOK, response.Success(map[string]string{"message": "status updated"}))
	}

	return c.JSON(http.StatusBadRequest, response.Error("status or assigned_to is required"))
}
