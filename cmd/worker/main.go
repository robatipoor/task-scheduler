package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/robatipoor/task-scheduler/internal/worker/client"
	"github.com/robatipoor/task-scheduler/internal/worker/config"
	"github.com/robatipoor/task-scheduler/internal/worker/controllers"
	"github.com/robatipoor/task-scheduler/internal/worker/models"
	"github.com/robatipoor/task-scheduler/internal/worker/repositories"
	"github.com/robatipoor/task-scheduler/internal/worker/routes"
	"github.com/robatipoor/task-scheduler/internal/worker/services"
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
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	masterClient := client.NewMasterClient(cfg)
	_, err = masterClient.Register(cfg.Master.Url, client.WorkerRegisterRequest{
		BaseUrl: fmt.Sprintf("%s%s", cfg.Server.Schema, addr),
	})
	if err != nil {
		log.Println("Failed to connect master node.")
	}

	taskRepo := repositories.NewTaskRepository(db)
	healthCheckService := services.NewHealthCheckService(db)
	taskService := services.NewTaskService(taskRepo)
	schedulerService := services.NewSchedulerService(ctx, cfg, masterClient, taskRepo)
	taskController := controllers.NewTaskController(taskService)
	healthCheckController := controllers.NewHealthCheckController(healthCheckService)

	log.Println("Setting up router...")
	r, err := routes.SetupRouter(healthCheckController, taskController)
	if err != nil {
		log.Fatalf("Error setuping routers: %v", err)
	}
	log.Println("Router setup complete.")

	schedulerService.Run()

	log.Printf("Starting server on %s:%s...\n", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
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

	if err := srv.ListenAndServe(); err != nil {
		log.Println("Error starting server: ", err)
	}
}
