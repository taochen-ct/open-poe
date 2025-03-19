package middleware

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"open-poe/internal/pkg/response"
)

type Recovery struct {
	loggerWriter *lumberjack.Logger
}

func NewRecoveryMiddleware(loggerWriter *lumberjack.Logger) *Recovery {
	return &Recovery{
		loggerWriter: loggerWriter,
	}
}

func (m *Recovery) Handler() gin.HandlerFunc {
	return gin.RecoveryWithWriter(
		m.loggerWriter,
		func(c *gin.Context, err interface{}) {
			c.JSON(http.StatusInternalServerError, response.Response{
				ErrorCode: response.ServerError,
				Data:      nil,
				Message:   "Internal Server Error",
			})
		},
	)
}
