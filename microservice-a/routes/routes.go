package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/worlder-team/microservice-server/microservice-a/configs"
	"github.com/worlder-team/microservice-server/microservice-a/docs"
	generatorHandlers "github.com/worlder-team/microservice-server/microservice-a/modules/generator/handlers"
	healthHandlers "github.com/worlder-team/microservice-server/microservice-a/modules/health/handlers"
)

// @title Microservice A API
// @version 1.0
// @description Sensor data generation service
// @BasePath /api/v1

// Router holds all dependencies needed for routing
type Router struct {
	generatorHandler *generatorHandlers.GeneratorHandler
	healthHandler    *healthHandlers.HealthHandler
	config           *configs.Config
}

// NewRouter creates a new router instance
func NewRouter(
	generatorHandler *generatorHandlers.GeneratorHandler,
	healthHandler *healthHandlers.HealthHandler,
	config *configs.Config,
) *Router {
	return &Router{
		generatorHandler: generatorHandler,
		healthHandler:    healthHandler,
		config:           config,
	}
}

// SetupRoutes configures all routes for the application
func (r *Router) SetupRoutes(e *echo.Echo) {
	// Update Swagger host dynamically
	docs.SwaggerInfo.Host = "localhost:" + r.config.Server.ExternalPort

	// Swagger documentation
	r.setupSwaggerRoutes(e)

	// Setup API versions
	r.setupV1Routes(e)
	// Future: r.setupV2Routes(e) - when you need API v2

	// Example for future API v2 (commented out):
	// r.setupV2Routes(e) - would have different handlers, maybe different response formats
}

// setupV1Routes configures API v1 routes
func (r *Router) setupV1Routes(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	// Module-specific routes for v1
	r.setupHealthRoutes(v1)
	r.setupGeneratorRoutes(v1)
}

// setupSwaggerRoutes configures Swagger documentation routes
func (r *Router) setupSwaggerRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

// setupHealthRoutes configures health check routes
func (r *Router) setupHealthRoutes(api *echo.Group) {
	api.GET("/health", r.healthHandler.Health)
}

// setupGeneratorRoutes configures generator routes
func (r *Router) setupGeneratorRoutes(api *echo.Group) {
	api.GET("/status", r.generatorHandler.GetStatus)
	api.POST("/frequency", r.generatorHandler.SetFrequency)
	api.GET("/frequency", r.generatorHandler.GetFrequency)
	api.POST("/start", r.generatorHandler.StartGeneration)
	api.POST("/stop", r.generatorHandler.StopGeneration)
}

// Example: Future API v2 implementation (commented out)
//
// // setupV2Routes configures API v2 routes
// func (r *Router) setupV2Routes(e *echo.Echo) {
// 	v2 := e.Group("/api/v2")
//
// 	// v2 might have different response formats, new features, etc.
// 	r.setupHealthRoutesV2(v2)
// 	r.setupGeneratorRoutesV2(v2)
// }
//
// // setupGeneratorRoutesV2 configures generator routes for API v2
// func (r *Router) setupGeneratorRoutesV2(api *echo.Group) {
// 	// v2 might have enhanced status with more details
// 	api.GET("/status", r.generatorHandler.GetStatusV2)
// 	// v2 might support batch frequency updates
// 	api.POST("/frequencies", r.generatorHandler.SetFrequenciesBatch)
// 	// etc.
// }
