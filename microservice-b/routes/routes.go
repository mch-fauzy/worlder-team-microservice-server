package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/worlder-team/microservice-server/microservice-b/configs"
	"github.com/worlder-team/microservice-server/microservice-b/docs"
	authHandlers "github.com/worlder-team/microservice-server/microservice-b/modules/auth/handlers"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/interfaces"
	healthHandlers "github.com/worlder-team/microservice-server/microservice-b/modules/health/handlers"
	sensorHandlers "github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/handlers"
	sharedMiddleware "github.com/worlder-team/microservice-server/shared/middleware"
)

// @title Microservice B API
// @version 1.0
// @description Sensor data storage and management service
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Router holds all dependencies needed for routing
type Router struct {
	sensorHandler *sensorHandlers.SensorHandler
	authHandler   *authHandlers.AuthHandler
	healthHandler *healthHandlers.HealthHandler
	jwtService    interfaces.JWTServiceInterface
	config        *configs.Config
}

// NewRouter creates a new router instance
func NewRouter(
	sensorHandler *sensorHandlers.SensorHandler,
	authHandler *authHandlers.AuthHandler,
	healthHandler *healthHandlers.HealthHandler,
	jwtService interfaces.JWTServiceInterface,
	config *configs.Config,
) *Router {
	return &Router{
		sensorHandler: sensorHandler,
		authHandler:   authHandler,
		healthHandler: healthHandler,
		jwtService:    jwtService,
		config:        config,
	}
}

// SetupRoutes configures all routes for the application
func (r *Router) SetupRoutes(e *echo.Echo) {
	// Update Swagger host dynamically
	docs.SwaggerInfo.Host = "localhost:" + r.config.Server.Port

	// Swagger documentation
	r.setupSwaggerRoutes(e)

	// Setup API versions
	r.setupV1Routes(e)
	// Future: r.setupV2Routes(e) - when you need API v2
}

// setupV1Routes configures API v1 routes
func (r *Router) setupV1Routes(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	// Module-specific routes for v1
	r.setupHealthRoutes(v1)
	r.setupAuthRoutes(v1)
	r.setupSensorRoutes(v1)
}

// setupSwaggerRoutes configures Swagger documentation routes
func (r *Router) setupSwaggerRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

// setupHealthRoutes configures health check routes
func (r *Router) setupHealthRoutes(api *echo.Group) {
	api.GET("/health", r.healthHandler.Health)
}

// setupAuthRoutes configures authentication routes
func (r *Router) setupAuthRoutes(api *echo.Group) {
	auth := api.Group("/auth")
	auth.POST("/login", r.authHandler.Login)
}

// setupSensorRoutes configures sensor data routes (protected)
func (r *Router) setupSensorRoutes(api *echo.Group) {
	sensors := api.Group("/sensors")
	sensors.Use(sharedMiddleware.JWTAuth(r.jwtService))

	sensors.GET("", r.sensorHandler.List)
	sensors.GET("/duration", r.sensorHandler.GetByDuration)
	sensors.GET("/:id", r.sensorHandler.GetByID)
	sensors.GET("/:id1/:id2", r.sensorHandler.GetByIDCombination)
	sensors.PATCH("/:id", r.sensorHandler.Update)
	sensors.DELETE("/:id", r.sensorHandler.Delete)
}
