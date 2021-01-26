package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

func SystemState(ctx *gin.Context){
	data := gin.H{}
	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// goroutine 的数量
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU
	ctx.JSON(http.StatusOK, data)
}
