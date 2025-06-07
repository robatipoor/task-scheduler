package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/robatipoor/task-scheduler/internal/worker/controllers"
)

func SetupRouter(healthCheckController *controllers.HealthCheckController, taskController *controllers.TaskController) (*gin.Engine, error) {
	router := gin.Default()
	router.GET("/health", healthCheckController.Check)
	api := router.Group("/api/v1")
	{
		api.POST("/tasks/submit", taskController.SubmitTask)
		api.GET("/tasks/result/:TrackID", taskController.GetResultTask)
	}

	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}

	return router, nil
}
