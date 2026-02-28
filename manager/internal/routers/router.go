package routers

import (
	"github.com/TKaterinna/CrackHash/manager/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewTaskRouter(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	router.POST("/api/hash/crack", taskHandler.Crack)
	router.GET("/api/hash/status", taskHandler.Status)
	router.PATCH("/internal/api/manager/hash/crack/request", taskHandler.UpdateResult)
}
