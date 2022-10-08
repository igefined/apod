package http

import (
	"github.com/igilgyrg/betera-test/internal/config"
	internalerror "github.com/igilgyrg/betera-test/internal/error"
	"github.com/igilgyrg/betera-test/internal/media"
	"github.com/igilgyrg/betera-test/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type mediaHandler struct {
	cfg            *config.Config
	usecase        model.MediaUsecase
	contextTimeout time.Duration
}

func NewMediaHandler(cfg *config.Config, mediaUC model.MediaUsecase, contextTimeout time.Duration) media.Handler {
	return &mediaHandler{cfg: cfg, usecase: mediaUC, contextTimeout: contextTimeout}
}

// APOD NASA godoc
// @Summary Get daily picture from NASA
// @Description Get daily picture from NASA
// @Tags nasa
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/apod [get]
func (m mediaHandler) APOD() echo.HandlerFunc {
	return func(c echo.Context) error {
		m, err := m.usecase.GetAPOD(c.Request().Context())
		if err != nil {
			return internalerror.NewSystemError(err)
		}

		return c.JSON(http.StatusOK, m)
	}
}
