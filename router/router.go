package router

import (
	"github.com/gin-gonic/gin"
	"go_chatserver/controller"
	"net/http"
)
func InitRouter() * gin.Engine{
	Router := gin.Default()
	Router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"success": true,
		})
	})
	systemRouter := Router.Group("/system")
	{
		systemRouter.GET("/state", controller.SystemState)
	}
	return Router
}
