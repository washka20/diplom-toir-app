package handlers

import (
	"context"
	"net/http"
	"strconv"

	"toir-app/internal/models"
	"toir-app/pkg/response"

	"github.com/labstack/echo/v4"
)

// EquipmentServiceInterface определяет контракт сервиса оборудования для handler.
type EquipmentServiceInterface interface {
	Create(ctx context.Context, eq *models.Equipment) error
	GetByID(ctx context.Context, id uint) (*models.Equipment, error)
	List(ctx context.Context, page, perPage int, status, location string) ([]models.Equipment, int64, error)
	Update(ctx context.Context, eq *models.Equipment) error
	Delete(ctx context.Context, id uint) error
}

// EquipmentHandler обрабатывает HTTP-запросы для оборудования.
type EquipmentHandler struct {
	service EquipmentServiceInterface
}

// NewEquipmentHandler создаёт handler оборудования.
func NewEquipmentHandler(service EquipmentServiceInterface) *EquipmentHandler {
	return &EquipmentHandler{service: service}
}

// Create обрабатывает POST /api/equipment.
func (h *EquipmentHandler) Create(c echo.Context) error {
	var eq models.Equipment
	if err := c.Bind(&eq); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if eq.Name == "" || eq.InventoryNumber == "" {
		return c.JSON(http.StatusBadRequest, response.Error("name and inventory_number are required"))
	}

	if err := h.service.Create(c.Request().Context(), &eq); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(eq))
}

// List обрабатывает GET /api/equipment.
func (h *EquipmentHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage < 1 {
		perPage = 20
	}

	status := c.QueryParam("status")
	location := c.QueryParam("location")

	items, total, err := h.service.List(c.Request().Context(), page, perPage, status, location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Paginated(items, page, perPage, total))
}

// GetByID обрабатывает GET /api/equipment/:id.
func (h *EquipmentHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid equipment id"))
	}

	eq, err := h.service.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(eq))
}

// Update обрабатывает PUT /api/equipment/:id.
func (h *EquipmentHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid equipment id"))
	}

	var eq models.Equipment
	if err := c.Bind(&eq); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	eq.ID = uint(id)

	if err := h.service.Update(c.Request().Context(), &eq); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(eq))
}

// Delete обрабатывает DELETE /api/equipment/:id.
func (h *EquipmentHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid equipment id"))
	}

	if err := h.service.Delete(c.Request().Context(), uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(nil))
}
