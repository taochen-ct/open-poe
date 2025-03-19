package user

import (
	"github.com/duke-git/lancet/v2/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"open-poe/config"
	"open-poe/internal/pkg/request"
	"open-poe/internal/pkg/response"
)

type Handler struct {
	logger  *zap.Logger
	conf    *config.Configuration
	service *Service
}

func NewHandler(logger *zap.Logger, conf *config.Configuration, service *Service) *Handler {
	return &Handler{
		logger:  logger,
		conf:    conf,
		service: service,
	}
}

func (h *Handler) Register(ctx *gin.Context) {
	var form request.Register
	if err := ctx.ShouldBindJSON(&form); err != nil {
		response.ServiceError(ctx, err)
		return
	}
	if validator.IsEmail(form.UserEmail) || validator.IsEmptyString(form.Username) {
		response.BadRequestError(
			ctx,
			response.BadRequest("params error; email or name maybe null", http.StatusBadRequest),
		)
		return
	}
	u, err := h.service.Register(ctx, &form)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, u)
	return
}
