package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

/* Look for api key, returns a 401 if one is not found */
func APIKeyMiddleware(apiKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestAPIKey := c.Request().Header.Get("X-API-Key")
			if requestAPIKey != apiKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API Key")
			}
			return next(c)
		}
	}
}
