package user

import (
	"context"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"open-poe/internal/pkg/request"
	"open-poe/internal/pkg/response"
)

type Repo interface {
	FindByUid(context.Context, string) (*User, error)
	FindById(context.Context, string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	Create(context.Context, *User) (*User, error)
}

type Service struct {
	repo   Repo
	rdb    *redis.Client
	sf     *sonyflake.Sonyflake
	logger *zap.Logger
}

func NewService(repo Repo, rdb *redis.Client, sf *sonyflake.Sonyflake, logger *zap.Logger) *Service {
	return &Service{repo: repo, rdb: rdb, sf: sf, logger: logger}
}

// Register create user
func (s *Service) Register(ctx *gin.Context, param *request.Register) (*User, error) {
	userInstance := &User{
		Name:     param.Username,
		Email:    param.UserEmail,
		Password: param.Password,
	}
	userId, err := s.sf.NextID()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, response.InternalServer("register fail; generate userid failed")
	}
	userInstance.Userid = convertor.ToString(userId)

	u, err := s.repo.Create(ctx, userInstance)
	if err != nil {
		return nil, response.BadRequest("register fail; create user failed; " + err.Error())

	}
	return u, nil
}
