package handlers

import (
	"context"
	"net/http"

	"toir-app/internal/models"
	"toir-app/pkg/response"

	"github.com/labstack/echo/v4"
)

// DashboardService определяет контракт сервиса дашборда.
type DashboardService interface {
	GetMetrics(ctx context.Context) (*models.DashboardMetrics, error)
}

// DashboardHandler обрабатывает HTTP-запросы дашборда.
type DashboardHandler struct {
	service DashboardService
}

// NewDashboardHandler создаёт handler дашборда.
func NewDashboardHandler(service DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetMetrics godoc
// @Summary Метрики дашборда
// @Description Получение агрегированных метрик для дашборда
// @Tags dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.DashboardMetrics}
// @Failure 500 {object} response.Response
// @Router /dashboard [get]
func (h *DashboardHandler) GetMetrics(c echo.Context) error {
	metrics, err := h.service.GetMetrics(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(metrics))
}
