package handler

import (
	"net/http"

	"github.com/abhinandpn/UnVocal/services/user-service/model"
	"github.com/abhinandpn/UnVocal/services/user-service/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "User registration details"
// @Success 201 {object} map[string]string "User created successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 409 {object} map[string]string "Email or phone number already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/new [post]
func (h *UserHandler) Register(c *gin.Context) {

	ctx := c.Request.Context()

	var user model.RegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Register(
		ctx,
		user.Name,
		user.Email,
		user.Password,
		user.Number,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

// GetUser godoc
// @Summary Get user by user code
// @Description Retrieves user details by user code.
// @Tags Users
// @Accept json
// @Produce json
// @Param uid path string true "User Code"
// @Success 200 {object} model.UserResponse "User details"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{uid} [get]
func (h *UserHandler) GetUser(c *gin.Context) {

	uid := c.Param("uid")
	ctx := c.Request.Context()

	user, err := h.service.GetUserByUserCode(ctx, uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Register godoc
// @Summary Delete a user
// @Description Deletes a user account by user code.
// @Tags Users
// @Accept json
// @Produce json
// @Param uid path string true "User Code"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{uid} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {

	uid := c.Param("uid")
	ctx := c.Request.Context()

	if err := h.service.DeleteUser(ctx, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Login godoc
// @Summary User login
// @Description Authenticates a user and returns a JWT token.
// @Tags Users
// @Accept json
// @Produce json
// @Param loginRequest body model.LoginRequest true "User login details"
// @Success 200 {object} model.LoginResponse "Login successful"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {

	ctx := c.Request.Context()

	var loginRequest model.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Login(ctx, loginRequest.Identifier, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Update User
func (h *UserHandler) UpdateUser(c *gin.Context) {

	id := c.Param("id")
	ctx := c.Request.Context()
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id

	if err := h.service.UpdateUser(ctx, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// UserProfile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile.
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} map[string]string
// @Router /users/profile [get]
func (h *UserHandler) UserProfile(c *gin.Context) {

	ctx := c.Request.Context()

	// Get user_code from JWT middleware
	userCode, exists := c.Get("user_code")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// Type assertion
	uc, ok := userCode.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user information",
		})
		return
	}

	// Get user profile
	user, err := h.service.UserProfile(ctx, uc)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

