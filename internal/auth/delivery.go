package auth

import "github.com/labstack/echo/v4"

type Handler interface {
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Register() echo.HandlerFunc
	Refresh() echo.HandlerFunc
}
