package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/igilgyrg/betera-test/internal/config"
	internalerror "github.com/igilgyrg/betera-test/internal/error"
	"github.com/igilgyrg/betera-test/pkg/logging"
	"github.com/labstack/echo/v4"
	"strings"
)

func (m *Manager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					logging.Log().Logger.Error("error of bearer token")
					return internalerror.NewUnauthorizedError(internalerror.Unauthorized)
				}

				tokenString := headerParts[1]

				err := m.validateJWTToken(tokenString, c, m.cfg)
				if err != nil {
					return internalerror.NewUnauthorizedError(err)
				}

				return next(c)
			}
			return internalerror.NewUnauthorizedError(internalerror.Unauthorized)
		}
	}
}

func (m *Manager) validateJWTToken(tokenString string, c echo.Context, cfg *config.Config) error {
	if tokenString == "" {
		return errors.New("")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.JWTSecretKey)
		return secret, nil
	})

	if err != nil {
		return errors.New("")
	}

	if !token.Valid {
		return errors.New("")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["sub"].(string)
		if !ok {
			return errors.New("")
		}

		user, err := m.userUC.Get(c.Request().Context(), uuid.MustParse(userID))
		if err != nil {
			return errors.New("")
		}

		c.Set("user", user)

		ctx := context.WithValue(c.Request().Context(), "user", user)
		c.SetRequest(c.Request().WithContext(ctx))
	}

	return nil
}
