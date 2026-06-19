package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"Clinic_backend/internal/service"
	"Clinic_backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo     repository.UserRepositoryInterface
	auditService service.AuditServiceInterface
}

func NewUserHandler(userRepo repository.UserRepositoryInterface, auditService service.AuditServiceInterface) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		auditService: auditService,
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
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
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
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID in token")
		return
	}

	var req entity.User
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedUser, err := h.userRepo.Update(c.Request.Context(), userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
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
	limit := 20
	offset := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	users, total, err := h.userRepo.GetAllPaginated(c.Request.Context(), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]entity.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *user.ToResponse()
	}

	utils.PaginatedResponse(c, http.StatusOK, responses, offset, limit, total)
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req entity.User
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedUser, err := h.userRepo.Update(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	adminID := c.GetInt("user_id")
	h.auditService.Log(c.Request.Context(), adminID, "UPDATE", "user", &id, c.ClientIP())
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userRepo.Delete(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	adminID := c.GetInt("user_id")
	h.auditService.Log(c.Request.Context(), adminID, "DELETE", "user", &id, c.ClientIP())
	c.Status(http.StatusNoContent)
}
