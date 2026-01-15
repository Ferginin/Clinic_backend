package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryServiceInterface
}

func NewCategoryHandler(categoryService service.CategoryServiceInterface) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create service category
// @Description Create a new service category (admin only)
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.ServiceCategory true "Category data"
// @Success 201 {object} entity.ServiceCategory
// @Failure 400 {object} map[string]string
// @Router /service-categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category entity.ServiceCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.categoryService.CreateCategory(c.Request.Context(), &category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetAllCategories godoc
// @Summary Get all service categories
// @Description Get list of all service categories
// @Tags categories
// @Produce json
// @Success 200 {array} entity.ServiceCategory
// @Router /service-categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get service category details by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} entity.ServiceCategory
// @Failure 404 {object} map[string]string
// @Router /service-categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetFavorites godoc
// @Summary Get favorite categories
// @Description Get list of favorite service categories
// @Tags categories
// @Produce json
// @Success 200 {array} entity.ServiceCategory
// @Router /service-categories/favorite [get]
func (h *CategoryHandler) GetFavorites(c *gin.Context) {
	categories, err := h.categoryService.GetFavoriteCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update service category (admin only)
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body entity.ServiceCategory true "Category update data"
// @Success 200 {object} entity.ServiceCategory
// @Failure 400 {object} map[string]string
// @Router /service-categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category entity.ServiceCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.categoryService.UpdateCategory(c.Request.Context(), id, &category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// ToggleFavorite godoc
// @Summary Toggle favorite status
// @Description Toggle favorite status of a category (admin only)
// @Tags categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /service-categories/{id}/favorite [patch]
func (h *CategoryHandler) ToggleFavorite(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.categoryService.ToggleFavorite(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite status toggled"})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete service category by ID (admin only)
// @Tags categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /service-categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
