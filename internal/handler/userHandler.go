package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserHandler(userRepo repository.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetMe godoc
// @Summary Get current user
// @Description Get currently authenticated user details
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} entity.UserResponse
// @Failure 401 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := h.userRepo.GetByID(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateMe godoc
// @Summary Update current user
// @Description Update currently authenticated user details
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.User true "User update data"
// @Success 200 {object} entity.UserResponse
// @Failure 400 {object} map[string]string
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req entity.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userRepo.Update(c.Request.Context(), userID.(int), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser.ToResponse())
}

// GetAll godoc
// @Summary Get all users
// @Description Get list of all users (admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.UserResponse
// @Failure 403 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]entity.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}

	c.JSON(http.StatusOK, responses)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get user details by ID (admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entity.UserResponse
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// Update godoc
// @Summary Update user
// @Description Update user by ID (admin only)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body entity.User true "User update data"
// @Success 200 {object} entity.UserResponse
// @Failure 400 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req entity.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userRepo.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser.ToResponse())
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by ID (admin only)
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
