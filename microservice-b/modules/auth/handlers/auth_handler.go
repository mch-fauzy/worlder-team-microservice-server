package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/dtos"
	"github.com/worlder-team/microservice-server/microservice-b/modules/auth/interfaces"
	"github.com/worlder-team/microservice-server/microservice-b/shared"
	"github.com/worlder-team/microservice-server/shared/constants"
)

type AuthHandler struct {
	authService interfaces.AuthServiceInterface
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService interfaces.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequest true "Login credentials (email and password)"
// @Success 200 {object} shared.APIResponse
// @Failure 400 {object} shared.APIResponse "Invalid request"
// @Failure 401 {object} shared.APIResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var request dtos.LoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrInvalidRequest,
			Error:   err.Error(),
		})
	}

	// Validate request
	if request.Email == "" || request.Password == "" {
		return c.JSON(http.StatusBadRequest, shared.APIResponse{
			Status:  constants.StatusError,
			Message: "Email and password are required",
		})
	}

	// Authenticate user
	response, err := h.authService.Login(c.Request().Context(), &request)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, shared.APIResponse{
			Status:  constants.StatusError,
			Message: constants.ErrUnauthorized,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.APIResponse{
		Status:  constants.StatusSuccess,
		Message: "Login successful",
		Data:    response,
	})
}
