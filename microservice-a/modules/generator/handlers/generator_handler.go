package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/dtos"
	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/interfaces"
	"github.com/worlder-team/microservice-server/microservice-a/shared"
	"github.com/worlder-team/microservice-server/shared/constants"
)

type GeneratorHandler struct {
	generatorService interfaces.GeneratorService
}

// NewGeneratorHandler creates a new generator handler
func NewGeneratorHandler(generatorService interfaces.GeneratorService) *GeneratorHandler {
	return &GeneratorHandler{
		generatorService: generatorService,
	}
}

// GetStatus godoc
// @Summary Get generation status
// @Description Get current generator status and configuration
// @Tags generator
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /status [get]
func (h *GeneratorHandler) GetStatus(c echo.Context) error {
	status := h.generatorService.GetStatus()

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Generator status retrieved successfully",
		Data:    status,
	})
}

// SetFrequency godoc
// @Summary Set generation frequency
// @Description Set the frequency of sensor data generation
// @Tags generator
// @Accept json
// @Produce json
// @Param request body dtos.FrequencyRequest true "Frequency parameters"
// @Success 200 {object} shared.APIResponse
// @Router /frequency [post]
func (h *GeneratorHandler) SetFrequency(c echo.Context) error {
	var request dtos.FrequencyRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   err.Error(),
		})
	}

	if err := h.generatorService.SetFrequency(request.Frequency); err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Frequency updated successfully",
		Data: dtos.FrequencyResponse{
			Frequency: request.Frequency,
			Duration:  h.generatorService.GetFrequency().String(),
		},
	})
}

// GetFrequency godoc
// @Summary Get current frequency
// @Description Get the current generation frequency
// @Tags generator
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /frequency [get]
func (h *GeneratorHandler) GetFrequency(c echo.Context) error {
	frequency := h.generatorService.GetFrequency()

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Current frequency retrieved successfully",
		Data: dtos.FrequencyResponse{
			Frequency: frequency.String(),
			Duration:  frequency.String(),
		},
	})
}

// StartGeneration godoc
// @Summary Start data generation
// @Description Start generating sensor data
// @Tags generator
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /start [post]
func (h *GeneratorHandler) StartGeneration(c echo.Context) error {
	if h.generatorService.IsRunning() {
		return c.JSON(http.StatusConflict, shared.APIResponse{
			Status:  constants.StatusError,
			Message: "Generator is already running",
		})
	}

	// Start generation in background with detached context
	go h.generatorService.StartGeneration(context.Background())

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Data generation started successfully",
	})
}

// StopGeneration godoc
// @Summary Stop data generation
// @Description Stop generating sensor data
// @Tags generator
// @Accept json
// @Produce json
// @Success 200 {object} shared.APIResponse
// @Router /stop [post]
func (h *GeneratorHandler) StopGeneration(c echo.Context) error {
	if !h.generatorService.IsRunning() {
		return c.JSON(http.StatusConflict, shared.APIResponse{
			Status:  constants.StatusError,
			Message: "Generator is not running",
		})
	}

	h.generatorService.StopGeneration()

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Data generation stopped successfully",
	})
}
