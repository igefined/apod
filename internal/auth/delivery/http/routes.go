package http

import (
	"github.com/igilgyrg/betera-test/internal/auth"
	"github.com/igilgyrg/betera-test/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroups *echo.Group, h auth.Handler, manager *middleware.Manager) {
	authGroups.POST("/login", h.Login(), manager.ResponseErrorToJSON())
	authGroups.POST("/register", h.Register(), manager.ResponseErrorToJSON())
	authGroups.POST("/logout", h.Logout(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	authGroups.POST("/refresh", h.Refresh(), manager.ResponseErrorToJSON())
}
