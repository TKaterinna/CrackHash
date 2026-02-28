package handlers

import (
	"net/http"

	"github.com/TKaterinna/CrackHash/worker/internal/models"
	"github.com/TKaterinna/CrackHash/worker/internal/services"
	"github.com/gin-gonic/gin"
)

type CalcHandler struct {
	service *services.CalcService
}

func NewCalcHandler(service *services.CalcService) *CalcHandler {
	return &CalcHandler{
		service: service,
	}
}

func (h *CalcHandler) Calc(ctx *gin.Context) {
	var req models.CrackTaskRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect data struct" + err.Error()})
		return
	}

	if err := h.service.Save(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
