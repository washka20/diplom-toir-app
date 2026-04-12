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

// Create godoc
// @Summary Создание оборудования
// @Description Создаёт новую единицу оборудования
// @Tags equipment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.Equipment true "Данные оборудования"
// @Success 201 {object} response.Response{data=models.Equipment}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /equipment [post]
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

// List godoc
// @Summary Список оборудования
// @Description Получение списка оборудования с пагинацией и фильтрами
// @Tags equipment
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param per_page query int false "Элементов на странице" default(20)
// @Param status query string false "Фильтр по статусу"
// @Param location query string false "Фильтр по расположению"
// @Success 200 {object} response.Response{data=[]models.Equipment,meta=response.Meta}
// @Failure 500 {object} response.Response
// @Router /equipment [get]
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

// GetByID godoc
// @Summary Получение оборудования по ID
// @Description Возвращает оборудование по идентификатору
// @Tags equipment
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID оборудования"
// @Success 200 {object} response.Response{data=models.Equipment}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /equipment/{id} [get]
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

// Update godoc
// @Summary Обновление оборудования
// @Description Обновляет данные единицы оборудования
// @Tags equipment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID оборудования"
// @Param request body models.Equipment true "Обновлённые данные"
// @Success 200 {object} response.Response{data=models.Equipment}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /equipment/{id} [put]
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

// Delete godoc
// @Summary Удаление оборудования
// @Description Удаляет единицу оборудования по ID
// @Tags equipment
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID оборудования"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /equipment/{id} [delete]
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
