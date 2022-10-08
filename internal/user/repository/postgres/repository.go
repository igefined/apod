package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/igilgyrg/betera-test/internal/model"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) model.UserRepository {
	return &postgresUserRepository{db: db}
}

func (p postgresUserRepository) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user *model.User
	res := p.db.WithContext(ctx).First(&user, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (p postgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	res := p.db.WithContext(ctx).Find(&user, "email = ?", email)
	if res.Error != nil || res.RowsAffected == 0 {
		return nil, res.Error
	}
	return user, nil
}

func (p postgresUserRepository) Update(ctx context.Context, user *model.User) error {
	return p.db.WithContext(ctx).Where("id = ?", user.ID).Updates(user).Error
}

func (p postgresUserRepository) Store(ctx context.Context, user *model.User) (uuid.UUID, error) {
	res := p.db.WithContext(ctx).Save(user)
	if res.Error != nil {
		return uuid.UUID{}, nil
	}
	return user.ID, nil
}

func (p postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return p.db.WithContext(ctx).Delete("id = ?", id).Error
}
