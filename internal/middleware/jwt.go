package middleware

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"open-poe/config"
	"open-poe/internal/cases/user"
	"open-poe/internal/compo"
	"open-poe/internal/pkg/response"
	"time"
)

type JWTAuthMiddleware struct {
	conf    *config.Configuration
	service *user.JwtService
	data    *compo.Data
}

func NewJWTAuthMiddleware(conf *config.Configuration, service *user.JwtService) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{conf: conf, service: service}
}

func (m *JWTAuthMiddleware) isInRefreshPeriod(exp *jwt.NumericDate, period int64) bool {
	return int64(exp.Sub(time.Now()).Seconds()) < period
}

func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.Request.Header.Get("Authorization")
		if validator.IsEmptyString(userToken) {
			response.BadRequestError(c, response.BadRequest("no authorization header", response.TokenError))
			return
		}
		token, err := jwt.ParseWithClaims(userToken, &user.CustomerClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.conf.Jwt.Secret), nil
		})
		claims := token.Claims.(*user.CustomerClaims)
		if err != nil || m.service.IsInBlackList(c, claims.ID, userToken) {
			response.BadRequestError(c, response.BadRequest("authorization expired", response.TokenError))
			return
		}

		// refresh
		if m.isInRefreshPeriod(claims.ExpiresAt, m.conf.Jwt.RefreshGracePeriod) {
			tokenData, err := m.service.RefreshToken(c, claims.ID, token)
			if err == nil {
				c.Header("new-token", tokenData.AccessToken)
				c.Header("new-expire-in", convertor.ToString(tokenData.ExpiresIn))
			}
		}
		c.Set("token", token)
		c.Set("id", claims.ID)
	}
}
