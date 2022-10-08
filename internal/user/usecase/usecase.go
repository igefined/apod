package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/igilgyrg/betera-test/internal/model"
	"time"
)

type userUsecase struct {
	userRepo       model.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepo model.UserRepository, timeout time.Duration) model.UserUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		contextTimeout: timeout,
	}
}

func (u userUsecase) Get(c context.Context, id uuid.UUID) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u userUsecase) GetByEmail(c context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u userUsecase) Update(c context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.userRepo.Update(ctx, user)
}

func (u userUsecase) Store(c context.Context, user *model.User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err := user.HashPassword()
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.userRepo.Store(ctx, user)
}

func (u userUsecase) Delete(c context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.userRepo.Delete(ctx, id)
}
