package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robatipoor/task-scheduler/internal/master/dto"
	internal_errors "github.com/robatipoor/task-scheduler/internal/master/errors"
	"github.com/robatipoor/task-scheduler/internal/master/services"
	"github.com/robatipoor/task-scheduler/internal/utils"
)

type WorkerController struct {
	workerService services.WorkerServiceInterface
}

func NewWorkerController(workerService services.WorkerServiceInterface) *WorkerController {
	return &WorkerController{workerService: workerService}
}

func (wc *WorkerController) Register(c *gin.Context) {
	var req dto.WorkerRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, internal_errors.RespondWithError(err.Error()))
		return
	}
	result, err := wc.workerService.Register(req)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, internal_errors.RespondWithError(utils.GetFirstItem(err.Error(), ":")))
		return
	}
	if result.CreatedAt == result.UpdatedAt {
		c.JSON(http.StatusCreated, dto.TaskResponse{ID: result.ID})
	} else {
		c.JSON(http.StatusOK, dto.TaskResponse{ID: result.ID})
	}

}

func (wc *WorkerController) UpdateTaskResult(c *gin.Context) {
	var req dto.WorkerUpdateResultTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, internal_errors.RespondWithError(err.Error()))
		return
	}
	result, err := wc.workerService.UpdateResultTask(req)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, internal_errors.RespondWithError(utils.GetFirstItem(err.Error(), ":")))
		return
	}
	c.JSON(http.StatusOK, dto.TaskResponse{ID: result.ID})
}
