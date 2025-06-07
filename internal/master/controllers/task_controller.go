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

type TaskController struct {
	taskService services.TaskServiceInterface
}

func NewTaskController(taskService services.TaskServiceInterface) *TaskController {
	return &TaskController{taskService: taskService}
}

func (tc *TaskController) SubmitTask(c *gin.Context) {
	var req dto.SubmitTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, internal_errors.RespondWithError(err.Error()))
		return
	}
	task, err := tc.taskService.CheckExistTask(req)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, internal_errors.RespondWithError(utils.GetFirstItem(err.Error(), ":")))
		return
	}
	log.Println(task)
	if task != nil {
		c.JSON(http.StatusOK, dto.TaskResponse{ID: task.ID})
		return
	}
	result, err := tc.taskService.SubmitTask(req)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, internal_errors.RespondWithError(utils.GetFirstItem(err.Error(), ":")))
		return
	}
	c.JSON(http.StatusCreated, dto.TaskResponse{ID: result.ID})
}

func (tc *TaskController) GetResultTask(c *gin.Context) {
	trackUID := c.Param("trackUID")
	result, err := tc.taskService.GetResultTask(trackUID)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
		c.JSON(http.StatusInternalServerError, internal_errors.RespondWithError(utils.GetFirstItem(err.Error(), ":")))
		return
	}
	if result != nil {
		c.JSON(http.StatusOK, dto.TaskResultResponse{Result: result.Result, Status: &result.Status})
	} else {
		c.JSON(http.StatusOK, dto.TaskResultResponse{})
	}

}
