package user

import (
	"github.com/duke-git/lancet/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"open-poe/config"
	"open-poe/internal/pkg/request"
	"open-poe/internal/pkg/request/user"
	"open-poe/internal/pkg/response"
	"time"
)

type Handler struct {
	logger  *zap.Logger
	conf    *config.Configuration
	service *Service
	auth    *JwtService
	rdb     *redis.Client
}

func NewHandler(
	logger *zap.Logger,
	conf *config.Configuration,
	service *Service,
	auth *JwtService,
	rdb *redis.Client,
) *Handler {
	return &Handler{
		logger:  logger,
		conf:    conf,
		service: service,
		auth:    auth,
		rdb:     rdb,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var form user.Register
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.ServiceError(ctx, err)
		return
	}
	if !validator.IsEmail(form.UserEmail) || validator.IsEmptyString(form.Username) {
		response.BadRequestError(
			ctx,
			response.BadRequest("params error; email or name maybe null", response.ValidateError),
		)
		return
	}
	u, err := h.service.Register(ctx, &form)
	if err != nil {
		response.BadRequestError(ctx, err)
		return
	}

	response.Success(ctx, u)
}

func (h *Handler) saveToken(c *gin.Context, userUid string, timeOut time.Duration, user *User) {
	pipe := h.rdb.Pipeline()
	pipe.HSet(
		c,
		userUid,
		map[string]interface{}{
			"username":   user.Name,
			"uid":        user.UID,
			"id":         user.Userid,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdateAt,
		},
	)
	pipe.Expire(c, userUid, timeOut)
	_, err := pipe.Exec(c)

	if err != nil {
		response.ServiceError(c, response.InternalServer("save token error; "+err.Error()))
		return
	} // execute in pipe
}

func (h *Handler) Login(ctx *gin.Context) {
	var form user.Login
	request.BindJson(ctx, &form, response.ServiceError)
	userInfo, err := h.service.Login(ctx, &form)
	if err != nil {
		response.BusinessError(
			ctx,
			err,
		)
		return
	}
	isLogin, err := h.rdb.Exists(ctx, userInfo.UID).Result()
	if err != nil {
		response.ServiceError(ctx, err)
		return
	}
	if isLogin == 1 {
		response.BusinessError(ctx, response.BusinessFail("repeat login"))
		return
	}
	tokenOutput, _, err := h.auth.CreateToken(userInfo)
	if err != nil {
		response.ServiceError(ctx, err)
		return
	}
	ctx.Header("Authorization", tokenOutput.AccessToken)
	h.saveToken(ctx, tokenOutput.Uid, time.Hour, userInfo)
	response.Success(ctx, tokenOutput)
}
