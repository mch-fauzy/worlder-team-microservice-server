package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/dtos"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/entities"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/interfaces"
	"github.com/worlder-team/microservice-server/microservice-b/shared"
	"github.com/worlder-team/microservice-server/shared/constants"
)

type SensorHandler struct {
	sensorService interfaces.SensorServiceInterface
}

// NewSensorHandler creates a new sensor handler
func NewSensorHandler(sensorService interfaces.SensorServiceInterface) *SensorHandler {
	return &SensorHandler{
		sensorService: sensorService,
	}
}

// List godoc
// @Summary List sensor data
// @Description List sensor data with pagination and filtering
// @Tags sensors
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param sensor_type query string false "Sensor type filter"
// @Param id1 query string false "ID1 filter"
// @Param id2 query int false "ID2 filter"
// @Param from_time query string false "From time filter (RFC3339)"
// @Param to_time query string false "To time filter (RFC3339)"
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} shared.APIResponse
// @Security Bearer
// @Router /sensors [get]
func (h *SensorHandler) List(c echo.Context) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")

	pagination := &dtos.PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	}

	// Parse filters
	filter := &dtos.SensorDataFilter{}

	if sensorType := c.QueryParam("sensor_type"); sensorType != "" {
		filter.SensorType = &sensorType
	}

	if id1 := c.QueryParam("id1"); id1 != "" {
		filter.ID1 = &id1
	}

	if id2Str := c.QueryParam("id2"); id2Str != "" {
		if id2, err := strconv.ParseInt(id2Str, 10, 32); err == nil {
			id2Int32 := int32(id2)
			filter.ID2 = &id2Int32
		}
	}

	if fromTimeStr := c.QueryParam("from_time"); fromTimeStr != "" {
		if fromTime, err := time.Parse(time.RFC3339, fromTimeStr); err == nil {
			filter.FromTime = &fromTime
		}
	}

	if toTimeStr := c.QueryParam("to_time"); toTimeStr != "" {
		if toTime, err := time.Parse(time.RFC3339, toTimeStr); err == nil {
			filter.ToTime = &toTime
		}
	}

	result, err := h.sensorService.ListSensorData(c.Request().Context(), filter, pagination)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInternalServer,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Sensor data retrieved successfully",
		Data:    result,
	})
}

// GetByIDCombination godoc
// @Summary Get sensor data by ID combination
// @Description Get sensor data by ID1 and ID2 combination
// @Tags sensors
// @Accept json
// @Produce json
// @Param id1 path string true "ID1"
// @Param id2 path int true "ID2"
// @Success 200 {object} shared.APIResponse
// @Security Bearer
// @Router /sensors/{id1}/{id2} [get]
func (h *SensorHandler) GetByIDCombination(c echo.Context) error {
	id1 := c.Param("id1")
	id2Str := c.Param("id2")

	id2, err := strconv.ParseInt(id2Str, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   "Invalid ID2 format",
		})
	}

	data, err := h.sensorService.GetSensorDataByIDCombination(c.Request().Context(), id1, int32(id2))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInternalServer,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Sensor data retrieved successfully",
		Data:    data,
	})
}

// GetByDuration godoc
// @Summary Get sensor data by duration
// @Description Get sensor data within a time range
// @Tags sensors
// @Accept json
// @Produce json
// @Param from_time query string true "From time (RFC3339)"
// @Param to_time query string true "To time (RFC3339)"
// @Success 200 {object} shared.APIResponse
// @Security Bearer
// @Router /sensors/duration [get]
func (h *SensorHandler) GetByDuration(c echo.Context) error {
	fromStr := c.QueryParam("from_time")
	toStr := c.QueryParam("to_time")

	if fromStr == "" || toStr == "" {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   "Both 'from_time' and 'to_time' parameters are required",
		})
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   "Invalid 'from' time format. Use RFC3339",
		})
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   "Invalid 'to' time format. Use RFC3339",
		})
	}

	data, err := h.sensorService.GetSensorDataByDuration(c.Request().Context(), from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInternalServer,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Sensor data retrieved successfully",
		Data:    data,
	})
}

// GetByID godoc
// @Summary Get sensor data by ID
// @Description Get single sensor data by ID
// @Tags sensors
// @Accept json
// @Produce json
// @Param id path int true "Sensor data ID"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse "Invalid ID"
// @Failure 404 {object} shared.APIResponse "Sensor data not found"
// @Security Bearer
// @Router /sensors/{id} [get]
func (h *SensorHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   "Invalid ID format",
		})
	}

	data, err := h.sensorService.GetSensorData(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, shared.APIResponse{
			Status:  constants.StatusError,
			Message: "Sensor data not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Sensor data retrieved successfully",
		Data:    data,
	})
}

// Update godoc
// @Summary Partially update sensor data
// @Description Partially update sensor data by ID (PATCH)
// @Tags sensors
// @Accept json
// @Produce json
// @Param id path int true "Sensor data ID"
// @Param request body UpdateRequest true "Update parameters (all fields optional)"
// @Success 200 {object} shared.APIResponse
// @Security Bearer
// @Router /sensors/{id} [patch]
func (h *SensorHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   "Invalid ID format",
		})
	}

	var request UpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   err.Error(),
		})
	}

	// Get existing record
	existingData, err := h.sensorService.GetSensorData(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, shared.APIResponse{
			Status:  constants.StatusError,
			Message: "Sensor data not found",
			Error:   err.Error(),
		})
	}

	// Create update data with only changed fields
	updateData := &entities.SensorData{
		SensorValue: existingData.SensorValue,
		SensorType:  existingData.SensorType,
		Timestamp:   existingData.Timestamp,
	}

	// Update only provided fields
	if request.SensorValue != nil {
		updateData.SensorValue = *request.SensorValue
	}
	if request.SensorType != nil {
		updateData.SensorType = *request.SensorType
	}
	if request.Timestamp != nil {
		updateData.Timestamp = *request.Timestamp
	}

	err = h.sensorService.UpdateSensorData(c.Request().Context(), uint(id), updateData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInternalServer,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Sensor data updated successfully",
	})
}

// Delete godoc
// @Summary Delete sensor data by ID
// @Description Soft delete sensor data by ID
// @Tags sensors
// @Accept json
// @Produce json
// @Param id path int true "Sensor data ID"
// @Success 200 {object} shared.APIResponse
// @Security Bearer
// @Router /sensors/{id} [delete]
func (h *SensorHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   "Invalid ID format",
		})
	}

	err = h.sensorService.DeleteSensorDataByID(c.Request().Context(), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInternalServer,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Sensor data deleted successfully",
		Data:    map[string]uint{"id": uint(id)},
	})
}

// Request/Response structures
type UpdateRequest struct {
	SensorValue *float64   `json:"sensor_value,omitempty"`
	SensorType  *string    `json:"sensor_type,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
}
