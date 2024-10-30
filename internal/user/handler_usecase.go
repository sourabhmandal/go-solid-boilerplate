package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userUseCase UserUseCase
}

func NewUserHandler(userUseCase UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// RegisterUser handles HTTP requests for registering a new user.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	// Define a struct to capture JSON data from the request body
	type RegisterUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Parse the request JSON into the RegisterUserRequest struct
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Call the RegisterUser method from the use case layer
	err := h.userUseCase.RegisterUser(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// GetUserByID handles HTTP requests to retrieve a user by their ID.
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Retrieve the userID from the path parameters
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Call the GetUserByID method from the use case layer
	user, err := h.userUseCase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user data in JSON format
	c.JSON(http.StatusOK, user)
}
