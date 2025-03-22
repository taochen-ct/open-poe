package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
	"open-poe/internal/cases/user"
	"open-poe/internal/middleware"
)

// ProviderSet is router providers.
var ProviderSet = wire.NewSet(CreateRouter)

func CreateRouter(
	recoveryMiddleware *middleware.Recovery,
	corsMiddleware *middleware.Cors,
	JWTAuthMiddleware *middleware.JWTAuthMiddleware,
	limiterMiddleware *middleware.Limiter,
	userHandler *user.Handler,
) *gin.Engine {

	// create new gin engine, api group
	router := gin.New()
	apiGroup := router.Group("/api/v1")
	// add middleware
	router.Use(
		gin.Logger(),                 // default logger
		recoveryMiddleware.Handler(), // logger
		corsMiddleware.Handler(),     // cors
		limiterMiddleware.Handler(),  // request rate limiter
	)
	// no auth
	// test service health
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	apiGroup.GET("/test", func(c *gin.Context) {
		res := make(map[string]string)
		res["message"] = "Hello World!"
		c.JSON(http.StatusOK, res)
	})
	apiGroup.POST("/user/register", userHandler.Register)
	apiGroup.POST("/user/login", userHandler.Login)

	// auth
	router.Use(JWTAuthMiddleware.Handler()) // auth middleware

	return router
}
