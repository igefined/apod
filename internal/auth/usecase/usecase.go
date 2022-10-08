package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/igilgyrg/betera-test/internal/auth"
	"github.com/igilgyrg/betera-test/internal/config"
	internalerror "github.com/igilgyrg/betera-test/internal/error"
	"github.com/igilgyrg/betera-test/internal/model"
	"github.com/igilgyrg/betera-test/pkg/logging"
	"time"
)

const (
	AccessTokenExpiredMinutes = 15 * time.Minute
	RefreshTokenExpiredHours  = 24 * time.Hour
)

type AuthUC struct {
	userRepository model.UserRepository
	cfg            *config.Config
	contextTimeout time.Duration
}

func NewAuthUseCase(userRepository model.UserRepository, cfg *config.Config, contextTimeout time.Duration) auth.UseCase {
	return &AuthUC{userRepository: userRepository, cfg: cfg, contextTimeout: contextTimeout}
}

func (a AuthUC) Login(ctx context.Context, user *model.User) (*model.Token, error) {
	u, err := a.userRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if err := u.ComparePassword(user.Password); err != nil {
		// TODO invalid password error
		return nil, err
	}

	accessToken, err := generateAccessToken(u.ID, a.cfg.JWTSecretKey)
	if err != nil {
		return nil, errors.New("")
	}

	refreshToken, err := generateRefreshToken(u, a.cfg.TokenSecretKey)
	if err != nil {
		return nil, errors.New("")
	}

	return &model.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a AuthUC) Register(ctx context.Context, user *model.User) (*model.Token, error) {
	err := user.HashPassword()
	if err != nil {
		return nil, err
	}

	id, err := a.userRepository.Store(ctx, user)
	user.ID = id
	if err != nil {
		return nil, err
	}

	accessToken, err := generateAccessToken(id, a.cfg.JWTSecretKey)
	if err != nil {
		return nil, internalerror.NewAccessTokenInvalid(err)
	}

	refreshToken, err := generateRefreshToken(user, a.cfg.TokenSecretKey)
	if err != nil {
		return nil, internalerror.NewRefreshTokenInvalid(err)
	}

	return &model.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a AuthUC) Refresh(c context.Context, refreshToken string) (*model.Token, error) {
	_, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	userID, ok := validateRefreshToken(refreshToken, a.cfg.TokenSecretKey)
	if ok {
		accessToken, err := generateAccessToken(uuid.MustParse(userID), a.cfg.JWTSecretKey)
		if err != nil {
			logging.Log().Warnf("error of generate access token")
			return nil, internalerror.NewAccessTokenInvalid(err)
		}

		return &model.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
	}
	return nil, internalerror.NewRefreshTokenInvalid(nil)
}

func generateAccessToken(userID uuid.UUID, accessTokenSignature string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(AccessTokenExpiredMinutes).Unix()
	claims["auth"] = true
	claims["sub"] = userID

	accessToken, err := token.SignedString([]byte(accessTokenSignature))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func generateRefreshToken(u *model.User, refreshTokenSignature string) (string, error) {
	rt := jwt.New(jwt.SigningMethodHS256)
	rtClaims := rt.Claims.(jwt.MapClaims)
	rtClaims["sub"] = u.ID
	rtClaims["exp"] = time.Now().Add(RefreshTokenExpiredHours).Unix()

	refreshToken, err := rt.SignedString([]byte(refreshTokenSignature))
	if err != nil {
		return "", err
	}

	return refreshToken, nil

}

func validateRefreshToken(refreshTokenString string, refreshTokenSignature string) (string, bool) {
	if refreshTokenString == "" {
		return "", false
	}

	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(refreshTokenSignature)
		return secret, nil
	})

	if err != nil {
		return "", false
	}

	if !refreshToken.Valid {
		return "", false
	}

	userID := ""

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		userID, ok = claims["sub"].(string)
		if !ok {
			return "", false
		}
	}

	return userID, true
}
