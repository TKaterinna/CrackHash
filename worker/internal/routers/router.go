package routers

import (
	"github.com/TKaterinna/CrackHash/worker/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewCalcRouter(router *gin.Engine, calcHandler *handlers.CalcHandler) {
	router.POST("/internal/api/worker/hash/crack/task", calcHandler.Calc)
}
