package middleware

import (
	"awesomeProject/internal/compo"
	"awesomeProject/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type Limiter struct {
	lm *compo.LimiterManager
}

func NewLimiterMiddleware(lm *compo.LimiterManager) *Limiter {
	return &Limiter{
		lm: lm,
	}
}

func (m *Limiter) Handler(key ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var limiterKey string
		if len(key) > 0 && len(key[0]) > 0 {
			limiterKey = key[0]
		} else {
			limiterKey = ctx.GetString("token")
			if len(limiterKey) == 0 {
				limiterKey = ctx.ClientIP()
			}
		}

		l := m.lm.GetLimiter(rate.Every(50*time.Millisecond), 300, limiterKey)

		if !l.L.Allow() {
			ctx.JSON(http.StatusTooManyRequests, response.Response{
				ErrorCode: response.TooManyRequests,
				Data:      nil,
				Message:   "Too many requests",
			})
			return
		}
	}
}
