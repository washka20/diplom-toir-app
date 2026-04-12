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

// GetMetrics обрабатывает GET /api/dashboard/metrics.
func (h *DashboardHandler) GetMetrics(c echo.Context) error {
	metrics, err := h.service.GetMetrics(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(metrics))
}
