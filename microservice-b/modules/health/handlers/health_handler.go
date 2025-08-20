package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/worlder-team/microservice-server/microservice-b/shared"
	"github.com/worlder-team/microservice-server/shared/constants"
)

type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary Health check
// @Description Health check endpoint
// @Tags health
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /health [get]
func (h *HealthHandler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"service":   "microservice-b",
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}
