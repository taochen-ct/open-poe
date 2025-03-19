package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterTestRoute(router *gin.RouterGroup) error {
	router.GET("/test", func(c *gin.Context) {
		res := make(map[string]string)
		res["message"] = "Hello World!"
		c.JSON(http.StatusOK, res)
	})
	return nil
}
