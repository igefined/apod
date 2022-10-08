package media

import "github.com/labstack/echo/v4"

type Handler interface {
	APOD() echo.HandlerFunc
}
