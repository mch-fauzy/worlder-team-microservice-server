package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/worlder-team/microservice-server/microservice-a/configs"
	generatorGrpc "github.com/worlder-team/microservice-server/microservice-a/modules/generator/grpc"
	generatorHandlers "github.com/worlder-team/microservice-server/microservice-a/modules/generator/handlers"
	generatorServices "github.com/worlder-team/microservice-server/microservice-a/modules/generator/services"
	healthHandlers "github.com/worlder-team/microservice-server/microservice-a/modules/health/handlers"
	"github.com/worlder-team/microservice-server/microservice-a/routes"
	sharedMiddleware "github.com/worlder-team/microservice-server/shared/middleware"
	"github.com/worlder-team/microservice-server/shared/utils"
)

func main() {
	// Load configuration
	cfg := configs.LoadConfig()

	// Initialize logger
	if err := utils.InitLogger(cfg.Server.LogLevel); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer utils.Sync()

	// Initialize gRPC client
	grpcClient, err := generatorGrpc.NewSensorClient(cfg.GetGRPCAddress())
	if err != nil {
		utils.Fatal("Failed to connect to gRPC server")
	}
	defer grpcClient.Close()

	// Initialize services
	generatorService := generatorServices.NewGeneratorService(grpcClient, cfg.Generator.SensorType, cfg.Generator.Frequency)

	// Initialize handlers
	generatorHandler := generatorHandlers.NewGeneratorHandler(generatorService)
	healthHandler := healthHandlers.NewHealthHandler()

	// Initialize router
	router := routes.NewRouter(generatorHandler, healthHandler, cfg)

	// Start data generation in background
	go generatorService.StartGeneration(context.Background())

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(sharedMiddleware.CORS())
	e.Use(sharedMiddleware.SecurityHeaders())
	e.Use(sharedMiddleware.RequestID())

	// Setup routes
	router.SetupRoutes(e)

	// Start HTTP server
	go func() {
		utils.Info(fmt.Sprintf("Starting HTTP server on port %s", cfg.Server.Port))
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Info("Shutting down server...")

	// Stop data generation
	generatorService.StopGeneration()

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		utils.Error("Server forced to shutdown")
	}

	utils.Info("Server stopped")
}
