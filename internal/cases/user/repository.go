package user

import (
	"context"
	"go.uber.org/zap"
	"open-poe/internal/compo"
)

type Repository struct {
	data *compo.Data
	log  *zap.Logger
}

func NewRepository(data *compo.Data, log *zap.Logger) Repo {
	return &Repository{
		data: data,
		log:  log,
	}
}

func (r Repository) Create(ctx context.Context, user *User) (*User, error) {
	if err := r.data.DB(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r Repository) FindByUid(ctx context.Context, uid string) (*User, error) {
	user := &User{}
	if err := r.data.DB(ctx).Model(user).Where("uid = ?", uid).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r Repository) FindById(ctx context.Context, id string) (*User, error) {
	user := &User{}
	if err := r.data.DB(ctx).Model(user).Where("user_id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r Repository) FindByUsername(ctx context.Context, name string) (*User, error) {
	user := &User{}
	if err := r.data.DB(ctx).Model(user).Where("user_name = ?", name).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
