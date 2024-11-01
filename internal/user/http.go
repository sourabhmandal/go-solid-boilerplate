package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type UserHandler interface {
	RegisterUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetAllUsers(c *gin.Context)

}



type userHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return &userHandler{userService: userService}
}

// RegisterUser handles HTTP requests for registering a new user.
func (h *userHandler) RegisterUser(c *gin.Context) {
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
	err := h.userService.RegisterUser(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Return a success response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// GetUserByID handles HTTP requests to retrieve a user by their ID.
func (h *userHandler) GetUserByID(c *gin.Context) {
	// Retrieve the userID from the path parameters
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Call the GetUserByID method from the use case layer
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user data in JSON format
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetAllUsers(c *gin.Context) {

	// Call the GetUserByID method from the use case layer
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}

	// Respond with the user data in JSON format
	c.JSON(http.StatusOK, users)
}