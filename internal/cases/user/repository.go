package user

import (
	"context"
	"github.com/duke-git/lancet/v2/cryptor"
	"go.uber.org/zap"
	"open-poe/internal/compo"
	"strconv"
	"time"
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

func (r *Repository) Create(ctx context.Context, user *User) (*User, error) {
	if err := r.data.DB(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindByUid(ctx context.Context, uid string) (*User, error) {
	user := &User{}
	if err := r.data.DB(ctx).Model(user).Where("uid = ?", uid).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindById(ctx context.Context, id string) (*User, error) {
	user := &User{}
	if err := r.data.DB(ctx).Model(user).Where("user_id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindByUsername(ctx context.Context, name string) (*User, error) {
	user := &User{}
	if err := r.data.DB(ctx).Model(user).Where("user_name = ?", name).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindByUserEmail(ctx context.Context, name string) (*User, error) {
	user := &User{}
	if err := r.data.DB(ctx).Model(user).Where("user_email = ?", name).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

type JwtRepository struct {
	data   *compo.Data
	logger *zap.Logger
}

func NewJwtRepository(data *compo.Data, logger *zap.Logger) JwtRepo {
	return &JwtRepository{
		data:   data,
		logger: logger,
	}
}

func (r *JwtRepository) getBlackListKey(userUid, token string) string {
	return "jwt_black_list_" + userUid + "_" + cryptor.Md5String(token)
}

func (r *JwtRepository) JoinBlackList(ctx context.Context, userUid, token string, joinUnix int64, expires time.Duration) error {
	return r.data.RDB().SetNX(ctx, r.getBlackListKey(userUid, token), joinUnix, expires).Err()
}

func (r *JwtRepository) GetBlackJoinUnix(ctx context.Context, userUid, token string) (int64, error) {
	joinUNixStr, err := r.data.RDB().Get(ctx, r.getBlackListKey(userUid, token)).Result()
	if err != nil {
		return 0, err
	}
	joinUnix, err := strconv.ParseInt(joinUNixStr, 10, 64)
	if err != nil || joinUnix == 0 {
		return 0, err
	}
	return joinUnix, nil
}
