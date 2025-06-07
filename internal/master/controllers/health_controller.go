package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robatipoor/task-scheduler/internal/master/dto"
	"github.com/robatipoor/task-scheduler/internal/master/services"
)

type HealthCheckController struct {
	healthCheckService *services.HealthCheckService
}

func NewHealthCheckController(healthCheckService *services.HealthCheckService) *HealthCheckController {
	return &HealthCheckController{
		healthCheckService: healthCheckService,
	}
}

func (hc *HealthCheckController) Check(c *gin.Context) {
	resp := hc.healthCheckService.Check()
	if anyServiceNotReady(resp.Services) {
		c.JSON(http.StatusInternalServerError, resp)
	} else {
		c.JSON(http.StatusOK, resp)
	}
}

func anyServiceNotReady(services []dto.ServiceCheckStatus) bool {
	for _, service := range services {
		if !service.IsReady {
			return true
		}
	}
	return false
}
