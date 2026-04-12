package middleware

import (
	"net/http"
	"strings"

	"toir-app/pkg/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTAuth возвращает middleware, которое проверяет JWT токен из заголовка Authorization.
// При успешной валидации устанавливает user_id (uint) и role (string) в echo.Context.
func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, response.Error("missing authorization header"))
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return c.JSON(http.StatusUnauthorized, response.Error("invalid authorization header format"))
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, response.Error("invalid or expired token"))
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.Error("invalid token claims"))
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.Error("invalid user_id in token"))
			}

			role, ok := claims["role"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, response.Error("invalid role in token"))
			}

			c.Set("user_id", uint(userIDFloat))
			c.Set("role", role)

			return next(c)
		}
	}
}
