package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/worlder-team/microservice-server/microservice-b/configs"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/entities"
	authHandlers "github.com/worlder-team/microservice-server/microservice-b/modules/auth/handlers"
	authServices "github.com/worlder-team/microservice-server/microservice-b/modules/auth/services"
	healthHandlers "github.com/worlder-team/microservice-server/microservice-b/modules/health/handlers"
	sensorEntities "github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/entities"
	sensorGrpc "github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/grpc"
	sensorHandlers "github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/handlers"
	sensorInterfaces "github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/interfaces"
	sensorRepositories "github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/repositories"
	sensorServices "github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/services"
	sharedServices "github.com/worlder-team/microservice-server/microservice-b/modules/shared/services"
	"github.com/worlder-team/microservice-server/microservice-b/routes"
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

	// Database connection
	db, err := initDatabase(cfg)
	if err != nil {
		utils.Fatal("Failed to connect to database")
	}

	// Run database seeding
	seeder := sharedServices.NewSeederService(db)
	if err := seeder.SeedAll(); err != nil {
		utils.Error(fmt.Sprintf("Database seeding failed: %v", err))
	}

	// Redis connection
	redisClient := initRedis(cfg)

	// Initialize repositories
	sensorRepo := sensorRepositories.NewSensorRepository(db)

	// Initialize JWT service
	jwtService := authServices.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer)

	// Initialize services
	sensorService := sensorServices.NewSensorService(sensorRepo)
	authService := authServices.NewAuthService(db, jwtService)

	// Initialize handlers
	sensorHandler := sensorHandlers.NewSensorHandler(sensorService)
	authHandler := authHandlers.NewAuthHandler(authService)
	healthHandler := healthHandlers.NewHealthHandler()

	// Initialize router
	router := routes.NewRouter(sensorHandler, authHandler, healthHandler, cfg)

	// Start gRPC server in goroutine
	go startGRPCServer(sensorService, cfg)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(sharedMiddleware.CORS())
	e.Use(sharedMiddleware.SecurityHeaders())
	e.Use(sharedMiddleware.RequestID())
	e.Use(sharedMiddleware.RateLimiter(redisClient, cfg.RateLimit.RequestsPerMinute))

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

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		utils.Error("Server forced to shutdown")
	}

	utils.Info("Server stopped")
}

func initDatabase(cfg *configs.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&sensorEntities.SensorData{}, &entities.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	utils.Info("Connected to database successfully")
	return db, nil
}

func initRedis(cfg *configs.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Redis.Password,
	})

	utils.Info("Connected to Redis successfully")
	return client
}

func startGRPCServer(sensorService sensorInterfaces.SensorServiceInterface, cfg *configs.Config) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Failed to listen on port %s", cfg.GRPC.Port))
	}

	s := grpc.NewServer()
	sensorGrpc.RegisterSensorServer(s, sensorService)

	utils.Info(fmt.Sprintf("Starting gRPC server on port %s", cfg.GRPC.Port))
	if err := s.Serve(lis); err != nil {
		utils.Fatal("Failed to serve gRPC")
	}
}
