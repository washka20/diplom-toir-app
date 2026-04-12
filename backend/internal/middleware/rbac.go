package middleware

import (
	"net/http"

	"toir-app/pkg/response"

	"github.com/labstack/echo/v4"
)

// RequireRole возвращает middleware, которое проверяет наличие роли пользователя
// в списке разрешённых ролей. Роль берётся из echo.Context (ключ "role").
func RequireRole(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok || role == "" {
				return c.JSON(http.StatusForbidden, response.Error("access denied: role not found"))
			}

			if _, exists := allowed[role]; !exists {
				return c.JSON(http.StatusForbidden, response.Error("access denied: insufficient permissions"))
			}

			return next(c)
		}
	}
}
