package handlers

import (
	"net/http"

	"github.com/TKaterinna/CrackHash/manager/internal/models"
	"github.com/TKaterinna/CrackHash/manager/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) Crack(ctx *gin.Context) {
	var req models.HashCrackRequest
	var err error
	var requestId uuid.UUID

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect data struct" + err.Error()})
		return
	}

	if requestId, err = h.service.Crack(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.HashCrackResponse{
		RequestId: requestId,
	})
}

func (h *TaskHandler) Status(ctx *gin.Context) {
	var requestId uuid.UUID
	var err error
	var status string
	var data []string
	var response models.HashStatusResponse

	strRequestId := ctx.Query("requestId")

	if len(strRequestId) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no query param"})
		return
	}

	if requestId, err = uuid.Parse(strRequestId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect query param"})
		return
	}

	if status, data, err = h.service.GetStatus(requestId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(data) == 0 {
		response = models.HashStatusResponse{
			Status: status,
		}
	} else {
		response = models.HashStatusResponse{
			Status: status,
			Data:   data,
		}
	}

	ctx.JSON(http.StatusOK, &response)
}

func (h *TaskHandler) UpdateResult(ctx *gin.Context) {
	var req models.CrackTaskResult

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect data struct" + err.Error()})
		return
	}

	if err := h.service.UpdateResult(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
