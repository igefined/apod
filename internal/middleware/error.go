package middleware

import (
	"encoding/json"
	"errors"
	internalerror "github.com/igilgyrg/betera-test/internal/error"
	"github.com/igilgyrg/betera-test/pkg/logging"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (m Manager) ResponseErrorToJSON() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var appError *internalerror.HttpError
			err := next(c)
			if err != nil {
				c.Response().Header().Set("Content-Type", "application/json")
				if errors.As(err, &appError) {
					c.Response().WriteHeader(appError.Status())
					c.Response().Write(appError.Marshal())
					logging.Log().Error(appError.Causes())
					return nil
				}

				c.Response().WriteHeader(http.StatusInternalServerError)
				errBytes, _ := json.Marshal(err)
				c.Response().Write(errBytes)
				logging.Log().Error(err)
			}
			c.Response().WriteHeader(http.StatusOK)
			return nil
		}
	}
}
