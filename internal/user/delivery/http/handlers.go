package http

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/igilgyrg/betera-test/internal/config"
	internalerror "github.com/igilgyrg/betera-test/internal/error"
	"github.com/igilgyrg/betera-test/internal/model"
	"github.com/igilgyrg/betera-test/internal/user"
	"github.com/igilgyrg/betera-test/pkg/utils"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"time"
)

type userHandler struct {
	cfg            *config.Config
	usecase        model.UserUsecase
	mediaUC        model.MediaUsecase
	contextTimeout time.Duration
}

func NewUserHandler(cfg *config.Config, usecase model.UserUsecase, mediaUC model.MediaUsecase, contextTimeout time.Duration) user.Handler {
	return &userHandler{cfg: cfg, usecase: usecase, mediaUC: mediaUC, contextTimeout: contextTimeout}
}

func UserFromContext(ctx echo.Context) (*model.User, error) {
	u := ctx.Get("user").(*model.User)
	if u == nil {
		return nil, errors.New("user have not founded from context")
	}
	return u, nil
}

// Current User godoc
// @Summary Get current authorized user
// @Description Get current authorized user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Router /api/v1/users/current [get]
func (u userHandler) Current() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := UserFromContext(c)
		if err != nil {
			return internalerror.NewUnauthorizedError(err)
		}

		userDB, err := u.usecase.Get(c.Request().Context(), user.ID)
		if err != nil {
			return internalerror.NewItemNotFound(err)
		}

		return c.JSON(http.StatusOK, userDB)
	}
}

// Albums User godoc
// @Summary Get current medias of authorized user
// @Description Get current medias of authorized user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} model.Media
// @Param date query string false "if empty, will response today; template 2006-01-02"
// @Router /api/v1/users/current/albums [get]
func (u userHandler) Albums() echo.HandlerFunc {
	return func(c echo.Context) error {
		userCtx, err := UserFromContext(c)
		if err != nil {
			return internalerror.NewUnauthorizedError(err)
		}

		var medias []model.Media

		dateAsString := c.QueryParam("date")
		if dateAsString != "" {
			date, err := time.Parse("2006-01-02", dateAsString)
			if err != nil {
				return internalerror.NewSystemError(err)
			}
			medias, err = u.mediaUC.ListByDate(c.Request().Context(), userCtx.ID, date)
			if err != nil {
				return internalerror.NewSystemError(err)
			}
		} else {
			medias, err = u.mediaUC.List(c.Request().Context(), userCtx.ID)
			if err != nil {
				return internalerror.NewSystemError(err)
			}
		}

		return c.JSON(http.StatusOK, medias)
	}
}

// DownloadImage User godoc
// @Summary Download image from local album
// @Description Download image from local album
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} []byte "array of bytes(image)"
// @Param date query string true "template 2006-01-02"
// @Param filename query string true "filename"
// @Router /api/v1/users/current/albums/download [get]
func (u userHandler) DownloadImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		dateAsString := c.QueryParam("date")
		filename := c.QueryParam("filename")
		userCtx, err := UserFromContext(c)
		if err != nil {
			return internalerror.NewUnauthorizedError(err)
		}

		date, err := time.Parse("2006-01-02", dateAsString)
		if err != nil {
			return internalerror.NewSystemError(err)
		}

		bytes, err := u.mediaUC.Download(c.Request().Context(), userCtx.ID, date, filename)
		if err != nil {
			return err
		}

		c.Response().Write(bytes)
		c.Response().WriteHeader(http.StatusOK)
		c.Response().Header().Set("Content-Type", "image/jpeg")
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		return nil
	}
}

// UploadImage User godoc
// @Summary Upload image to user current album
// @Description Upload image to user current album
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Param date formData string true "template 2006-01-02"
// @Param file formData file true "file"
// @Router /api/v1/users/current/albums/upload [post]
func (u userHandler) UploadImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		dateAsString := c.FormValue("date")
		multipartFile, err := c.FormFile("file")
		if err != nil {
			return internalerror.NewBadRequestError(err)
		}

		userCtx, err := UserFromContext(c)
		if err != nil {
			return internalerror.NewUnauthorizedError(err)
		}

		date, err := time.Parse("2006-01-02", dateAsString)
		if err != nil {
			return internalerror.NewSystemError(err)
		}

		file, err := multipartFile.Open()

		byteContent, err := ioutil.ReadAll(file)

		err = u.mediaUC.Store(c.Request().Context(), userCtx.ID, date, multipartFile.Filename, byteContent)
		if err != nil {
			return internalerror.NewSystemError(err)
		}

		return c.NoContent(http.StatusOK)
	}
}

// Get User godoc
// @Summary Get user by id
// @Description Get user by id
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Param id path string true "ID of user"
// @Router /api/v1/users [get]
func (u userHandler) Get() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")

		userDB, err := u.usecase.Get(c.Request().Context(), uuid.MustParse(userID))
		if err != nil {
			return internalerror.NewItemNotFound(err)
		}

		return c.JSON(http.StatusOK, userDB)
	}
}

// GetByEmail User godoc
// @Summary Get user by email
// @Description Get user by email
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Param email query string true "email of user"
// @Router /api/v1/users [get]
func (u userHandler) GetByEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.QueryParam("email")

		userDB, err := u.usecase.GetByEmail(c.Request().Context(), email)
		if err != nil {
			return internalerror.NewItemNotFound(err)
		}

		return c.JSON(http.StatusOK, userDB)
	}
}

// Save User godoc
// @Summary save user
// @Description save user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Param data body CreateUser true "Event repeat"
// @Router /api/v1/users [post]
func (u userHandler) Save() echo.HandlerFunc {
	return func(c echo.Context) error {
		userPatch := &CreateUser{}
		if err := utils.ReadRequest(c, userPatch); err != nil {
			return internalerror.NewBadRequestError(err)
		}

		id, err := u.usecase.Store(c.Request().Context(), &model.User{
			ID:        uuid.New(),
			FirstName: userPatch.FirstName,
			LastName:  userPatch.LastName,
			Email:     userPatch.Email,
			Password:  userPatch.Password,
			CreatedAt: time.Now(),
		})

		if err != nil {
			return internalerror.NewItemNotFound(err)
		}

		return c.JSON(http.StatusOK, id)
	}
}

// Update User godoc
// @Summary Update user
// @Description update user
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Param data body UpdateUser true "update user"
// @Router /api/v1/users [put]
func (u userHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		userCtx, err := UserFromContext(c)
		if err != nil {
			return internalerror.NewUnauthorizedError(err)
		}
		userPatch := &UpdateUser{}
		if err := utils.ReadRequest(c, userPatch); err != nil {
			return internalerror.NewBadRequestError(err)
		}

		userCtx.FirstName = userPatch.FirstName
		userCtx.LastName = userPatch.LastName
		userCtx.Email = userPatch.Email

		err = u.usecase.Update(c.Request().Context(), userCtx)
		if err != nil {
			return internalerror.NewSystemError(err)
		}

		return c.NoContent(http.StatusOK)
	}
}

// Delete User godoc
// @Summary Delete user
// @Description update user
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Param id path string true "delete user"
// @Router /api/v1/users [delete]
func (u userHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		err := u.usecase.Delete(c.Request().Context(), uuid.MustParse(id))
		if err != nil {
			return internalerror.NewSystemError(err)
		}

		return c.NoContent(http.StatusOK)
	}
}
