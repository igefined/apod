package http

import (
	"github.com/google/uuid"
	"github.com/igilgyrg/betera-test/internal/auth"
	"github.com/igilgyrg/betera-test/internal/config"
	internalerror "github.com/igilgyrg/betera-test/internal/error"
	"github.com/igilgyrg/betera-test/internal/model"
	"github.com/igilgyrg/betera-test/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type authHandler struct {
	cfg         *config.Config
	usecase     auth.UseCase
	userUsecase model.UserUsecase
}

func NewAuthHandler(cfg *config.Config, usecase auth.UseCase, userUsecase model.UserUsecase) auth.Handler {
	return &authHandler{cfg: cfg, usecase: usecase, userUsecase: userUsecase}
}

// Login auth godoc
// @Summary user login
// @Description login user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Token
// @Param data body Login true "login user"
// @Router /api/v1/auth/login [post]
func (a authHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			return internalerror.NewBadRequestError(err)
		}

		token, err := a.usecase.Login(c.Request().Context(), &model.User{
			Email:    login.Email,
			Password: login.Password,
		})

		if err != nil {
			return internalerror.NewUnauthorizedError(err)
		}

		return c.JSON(http.StatusOK, token)
	}
}

// Logout logout godoc
// @Summary user logout
// @Description logout user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/auth/logout [post]
func (a authHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
}

// Register registration godoc
// @Summary user registration
// @Description registation of user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Token
// @Param data body UserRegistration true "login user"
// @Router /api/v1/auth/register [post]
func (a authHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &UserRegistration{}
		if err := utils.ReadRequest(c, user); err != nil {
			return internalerror.NewBadRequestError(err)
		}

		userByEmail, err := a.userUsecase.GetByEmail(c.Request().Context(), user.Email)
		if err != nil {
			return internalerror.NewSystemError(err)
		}

		if userByEmail != nil {
			return internalerror.NewUserIsExistsWithEmail(err)
		}

		token, err := a.usecase.Register(c.Request().Context(), &model.User{
			ID:        uuid.New(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: time.Now(),
		})

		if err != nil {
			return internalerror.NewSystemError(err)
		}

		return c.JSON(http.StatusOK, token)
	}
}

// Refresh refresh godoc
// @Summary user refresh token
// @Description refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Token
// @Param data body Refresh true "refresh token"
// @Router /api/v1/auth/refresh [post]
func (a authHandler) Refresh() echo.HandlerFunc {
	return func(c echo.Context) error {
		refresh := &Refresh{}
		if err := utils.ReadRequest(c, refresh); err != nil {
			return internalerror.NewBadRequestError(err)
		}

		token, err := a.usecase.Refresh(c.Request().Context(), refresh.RefreshToken)
		if err != nil {
			return internalerror.NewUnauthorizedError(err)
		}

		return c.JSON(http.StatusOK, token)
	}
}
