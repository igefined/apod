package app

import (
	authHttp "github.com/igilgyrg/betera-test/internal/auth/delivery/http"
	authUseCase "github.com/igilgyrg/betera-test/internal/auth/usecase"
	"github.com/igilgyrg/betera-test/internal/config"
	mediaHttp "github.com/igilgyrg/betera-test/internal/media/delivery/http"
	mediaUsecase "github.com/igilgyrg/betera-test/internal/media/usecase"
	"github.com/igilgyrg/betera-test/internal/middleware"
	userHttp "github.com/igilgyrg/betera-test/internal/user/delivery/http"
	userRepository "github.com/igilgyrg/betera-test/internal/user/repository/postgres"
	userUsecase "github.com/igilgyrg/betera-test/internal/user/usecase"
	"github.com/labstack/echo/v4"
	"time"
)

func (a *App) MapHandlers(e *echo.Echo, cfg *config.Config, ctxTimeout time.Duration) error {
	// Init repositories
	userRepo := userRepository.NewPostgresUserRepository(a.postgresDB)

	//Init usecases
	authCase := authUseCase.NewAuthUseCase(userRepo, cfg, ctxTimeout)
	userCase := userUsecase.NewUserUsecase(userRepo, ctxTimeout)
	mediaCase := mediaUsecase.NewMediaUsecase(a.nasaClient, a.s3Storage, ctxTimeout)

	// Init handlers
	authHandlers := authHttp.NewAuthHandler(a.cfg, authCase, userCase)
	userHandlers := userHttp.NewUserHandler(a.cfg, userCase, mediaCase, ctxTimeout)
	mediaHandlers := mediaHttp.NewMediaHandler(a.cfg, mediaCase, ctxTimeout)

	mw := middleware.NewMiddlewareManager([]string{"*"}, a.cfg, userCase)

	v1 := e.Group("/api/v1")

	authGroup := v1.Group("/auth")
	userGroup := v1.Group("/users")
	mediaGroup := v1.Group("/media")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)
	userHttp.MapUserRoutes(userGroup, userHandlers, mw)
	mediaHttp.MapMediaRoutes(mediaGroup, mediaHandlers, mw)

	return nil
}
