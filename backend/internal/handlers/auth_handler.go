package handlers

import (
	"net/http"
	"strings"

	"toir-app/internal/services"
	"toir-app/pkg/response"

	"github.com/labstack/echo/v4"
)

// AuthHandler обрабатывает HTTP-запросы аутентификации.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler создаёт handler аутентификации.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginRequest представляет тело запроса на вход.
type LoginRequest struct {
	Username string `json:"username" example:"engineer"`
	Password string `json:"password" example:"password123"`
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Получение JWT токенов по логину и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} response.Response{data=services.TokenPair}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, response.Error("username and password are required"))
	}

	tokens, err := h.authService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(tokens))
}
