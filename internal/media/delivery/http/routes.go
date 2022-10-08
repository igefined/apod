package http

import (
	"github.com/igilgyrg/betera-test/internal/media"
	"github.com/igilgyrg/betera-test/internal/middleware"
	"github.com/labstack/echo/v4"
)

func MapMediaRoutes(groups *echo.Group, h media.Handler, manager *middleware.Manager) {
	groups.GET("/apod", h.APOD(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
}
