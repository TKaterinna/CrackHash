package routers

import (
	"github.com/TKaterinna/CrackHash/manager/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewTaskRouter(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	router.POST("/api/hash/crack", taskHandler.Crack)
	router.GET("/api/hash/status", taskHandler.Status)
	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})
}
