package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/robatipoor/task-scheduler/internal/master/client"
	"github.com/robatipoor/task-scheduler/internal/master/config"
	"github.com/robatipoor/task-scheduler/internal/master/controllers"
	"github.com/robatipoor/task-scheduler/internal/master/models"
	"github.com/robatipoor/task-scheduler/internal/master/repositories"
	"github.com/robatipoor/task-scheduler/internal/master/routes"
	"github.com/robatipoor/task-scheduler/internal/master/services"
)

func main() {
	log.Println("Starting application...")

	log.Println("Loading environment configuration...")
	cfg, err := config.LoadConfigure(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Println("Configuration loaded successfully.")

	log.Println("Setting up database connection...")
	db, err := models.SetupDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connection established successfully.")

	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	workerClient := client.NewWorkerClient(cfg)
	workerRepo := repositories.NewWorkerRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	assingTaskRepo := repositories.NewAssingTaskRepository(db)
	healthCheckService := services.NewHealthCheckService(db)
	taskService := services.NewTaskService(taskRepo)
	workerService := services.NewWorkerService(workerRepo, assingTaskRepo)
	schedulerService := services.NewSchedulerService(ctx, cfg, workerClient, workerRepo, taskRepo, assingTaskRepo)
	taskController := controllers.NewTaskController(taskService)
	workerController := controllers.NewWorkerController(workerService)
	healthCheckController := controllers.NewHealthCheckController(healthCheckService)

	log.Println("Setting up router...")
	r, err := routes.SetupRouter(healthCheckController, taskController, workerController)
	if err != nil {
		log.Fatalf("Error setuping routers: %v", err)
	}
	log.Println("Router setup complete.")

	schedulerService.Run()

	log.Printf("Starting server on %s:%s...\n", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: r,
	}
	go func() {
		<-quit
		log.Printf("Shutting down server...")
		cancel()
		schedulerService.Wait()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("Server forced to shutdown: ", err)
		}
		log.Printf("Server exiting")
	}()

	if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Println("Error starting server: ", err)
	}
}
