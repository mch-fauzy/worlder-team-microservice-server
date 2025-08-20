package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/worlder-team/microservice-server/microservice-a/shared"
	"github.com/worlder-team/microservice-server/shared/constants"
)

type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /health [get]
func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"service":   "microservice-a",
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}
