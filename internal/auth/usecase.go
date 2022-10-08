package auth

import (
	"context"
	"github.com/igilgyrg/betera-test/internal/model"
)

type UseCase interface {
	Login(ctx context.Context, user *model.User) (*model.Token, error)
	Register(ctx context.Context, user *model.User) (*model.Token, error)
	Refresh(ctx context.Context, refreshToken string) (*model.Token, error)
}
