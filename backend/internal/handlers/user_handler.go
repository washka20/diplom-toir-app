package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"toir-app/internal/models"
	"toir-app/internal/repository"
	"toir-app/pkg/response"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler обрабатывает HTTP-запросы управления пользователями (admin).
type UserHandler struct {
	repo repository.UserRepository
}

// NewUserHandler создаёт handler пользователей.
func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// CreateUserInput представляет тело запроса на создание пользователя.
type CreateUserInput struct {
	Username string `json:"username" example:"ivanov"`
	Email    string `json:"email" example:"ivanov@example.com"`
	Password string `json:"password" example:"securepass123"`
	FullName string `json:"full_name" example:"Иванов Иван Иванович"`
	Role     string `json:"role" example:"technician"`
}

// UpdateUserInput представляет тело запроса на обновление пользователя.
type UpdateUserInput struct {
	Email    string `json:"email" example:"new@example.com"`
	FullName string `json:"full_name" example:"Иванов Иван Иванович"`
	Role     string `json:"role" example:"engineer"`
	IsActive *bool  `json:"is_active" example:"true"`
	Password string `json:"password" example:"newpassword"`
}

// List godoc
// @Summary Список пользователей
// @Description Получение списка всех пользователей (только admin)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]models.User}
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *UserHandler) List(c echo.Context) error {
	users, err := h.repo.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(users))
}

// Create godoc
// @Summary Создание пользователя
// @Description Создаёт нового пользователя (только admin)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserInput true "Данные пользователя"
// @Success 201 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [post]
func (h *UserHandler) Create(c echo.Context) error {
	var input CreateUserInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(input.Email)
	input.FullName = strings.TrimSpace(input.FullName)
	input.Role = strings.TrimSpace(input.Role)

	if input.Username == "" || input.Email == "" || input.Password == "" || input.FullName == "" || input.Role == "" {
		return c.JSON(http.StatusBadRequest, response.Error("username, email, password, full_name and role are required"))
	}

	validRoles := map[string]struct{}{
		"admin": {}, "engineer": {}, "technician": {}, "operator": {},
	}
	if _, ok := validRoles[input.Role]; !ok {
		return c.JSON(http.StatusBadRequest, response.Error("role must be one of: admin, engineer, technician, operator"))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error("failed to hash password"))
	}

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hash),
		FullName:     input.FullName,
		Role:         input.Role,
		IsActive:     true,
	}

	if err := h.repo.Create(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(user))
}

// Update godoc
// @Summary Обновление пользователя
// @Description Обновляет данные пользователя (только admin)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Param request body UpdateUserInput true "Обновлённые данные"
// @Success 200 {object} response.Response{data=models.User}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid user id"))
	}

	user, err := h.repo.FindByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Error("user not found"))
	}

	var input UpdateUserInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request body"))
	}

	if email := strings.TrimSpace(input.Email); email != "" {
		user.Email = email
	}
	if fullName := strings.TrimSpace(input.FullName); fullName != "" {
		user.FullName = fullName
	}
	if role := strings.TrimSpace(input.Role); role != "" {
		validRoles := map[string]struct{}{
			"admin": {}, "engineer": {}, "technician": {}, "operator": {},
		}
		if _, ok := validRoles[role]; !ok {
			return c.JSON(http.StatusBadRequest, response.Error("role must be one of: admin, engineer, technician, operator"))
		}
		user.Role = role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}
	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.Error("failed to hash password"))
		}
		user.PasswordHash = string(hash)
	}

	if err := h.repo.Update(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(user))
}
