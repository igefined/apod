package http

import (
	"github.com/igilgyrg/betera-test/internal/middleware"
	"github.com/igilgyrg/betera-test/internal/user"
	"github.com/labstack/echo/v4"
)

func MapUserRoutes(userGroups *echo.Group, h user.Handler, manager *middleware.Manager) {
	userGroups.GET("/current", h.Current(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.GET("/current/albums", h.Albums(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.GET("/current/albums/download", h.DownloadImage(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.POST("/current/albums/upload", h.UploadImage(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.GET("/:id", h.Get(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.GET("", h.GetByEmail(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.PUT("/:id", h.Update(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.DELETE("/:id", h.Delete(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
	userGroups.POST("", h.Save(), manager.ResponseErrorToJSON(), manager.AuthJWTMiddleware())
}
