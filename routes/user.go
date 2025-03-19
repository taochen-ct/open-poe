package routes

import (
	"github.com/gin-gonic/gin"
	"open-poe/internal/cases/user"
)

func RegisterUserRoute(router *gin.RouterGroup, handler *user.Handler) {
	router.POST("/user/register", handler.Register)
	return
}
