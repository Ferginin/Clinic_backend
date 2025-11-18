package handler

import (
	"Clinic_backend/config"
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SpecializationHandler struct {
	specService *service.SpecializationService
}

func NewSpecializationHandler(cfg *config.Config, db *pgxpool.Pool) *SpecializationHandler {
	return &SpecializationHandler{
		specService: service.NewSpecializationService(db),
	}
}

// CreateSpecialization godoc
// @Summary Create specialization
// @Description Create a new specialization (admin only)
// @Tags specializations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.Specialization true "Specialization data"
// @Success 201 {object} entity.Specialization
// @Failure 400 {object} map[string]string
// @Router /specializations [post]
func (h *SpecializationHandler) CreateSpecialization(c *gin.Context) {
	var spec entity.Specialization
	if err := c.ShouldBindJSON(&spec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.specService.CreateSpecialization(c.Request.Context(), &spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetAllSpecializations godoc
// @Summary Get all specializations
// @Description Get list of all specializations
// @Tags specializations
// @Produce json
// @Success 200 {array} entity.Specialization
// @Router /specializations [get]
func (h *SpecializationHandler) GetAllSpecializations(c *gin.Context) {
	specializations, err := h.specService.GetAllSpecializations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, specializations)
}

// GetSpecializationByID godoc
// @Summary Get specialization by ID
// @Description Get specialization details by ID
// @Tags specializations
// @Produce json
// @Param id path int true "Specialization ID"
// @Success 200 {object} entity.Specialization
// @Failure 404 {object} map[string]string
// @Router /specializations/{id} [get]
func (h *SpecializationHandler) GetSpecializationByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid specialization ID"})
		return
	}

	spec, err := h.specService.GetSpecializationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Specialization not found"})
		return
	}

	c.JSON(http.StatusOK, spec)
}

// UpdateSpecialization godoc
// @Summary Update specialization
// @Description Update specialization (admin only)
// @Tags specializations
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Specialization ID"
// @Param request body entity.Specialization true "Specialization update data"
// @Success 200 {object} entity.Specialization
// @Failure 400 {object} map[string]string
// @Router /specializations/{id} [put]
func (h *SpecializationHandler) UpdateSpecialization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid specialization ID"})
		return
	}

	var spec entity.Specialization
	if err := c.ShouldBindJSON(&spec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.specService.UpdateSpecialization(c.Request.Context(), id, &spec)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteSpecialization godoc
// @Summary Delete specialization
// @Description Delete specialization by ID (admin only)
// @Tags specializations
// @Security BearerAuth
// @Param id path int true "Specialization ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /specializations/{id} [delete]
func (h *SpecializationHandler) DeleteSpecialization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid specialization ID"})
		return
	}

	if err := h.specService.DeleteSpecialization(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
