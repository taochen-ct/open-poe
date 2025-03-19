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
	limiterMiddleware *middleware.Limiter,
	userHandler *user.Handler,
) *gin.Engine {

	// create new gin engine
	router := gin.New()
	// add middleware
	router.Use(
		gin.Logger(),                 // default logger
		recoveryMiddleware.Handler(), // logger
		corsMiddleware.Handler(),     // cors
		limiterMiddleware.Handler(),  // request rate limiter
	)
	// test service health
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	// create api group
	apiGroup := router.Group("/api/v1")
	RegisterTestRoute(apiGroup)
	RegisterUserRoute(apiGroup, userHandler)
	return router
}
