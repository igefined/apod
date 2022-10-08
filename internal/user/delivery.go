package user

import "github.com/labstack/echo/v4"

type Handler interface {
	Get() echo.HandlerFunc
	Current() echo.HandlerFunc
	Albums() echo.HandlerFunc
	DownloadImage() echo.HandlerFunc
	UploadImage() echo.HandlerFunc
	GetByEmail() echo.HandlerFunc
	Save() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}
