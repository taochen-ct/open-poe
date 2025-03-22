package request

import "github.com/gin-gonic/gin"

func BindJson(ctx *gin.Context, v interface{}, fn func(*gin.Context, error)) {
	if err := ctx.ShouldBindJSON(v); err != nil {
		fn(ctx, err)
	}
}
