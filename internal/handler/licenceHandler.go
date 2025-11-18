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

type LicenseHandler struct {
	licenseService *service.LicenseService
}

func NewLicenseHandler(cfg *config.Config, db *pgxpool.Pool) *LicenseHandler {
	return &LicenseHandler{
		licenseService: service.NewLicenseService(db),
	}
}

// CreateLicense godoc
// @Summary Create license
// @Description Create a new clinic license (admin only)
// @Tags licenses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.License true "License data"
// @Success 201 {object} entity.License
// @Failure 400 {object} map[string]string
// @Router /licenses [post]
func (h *LicenseHandler) CreateLicense(c *gin.Context) {
	var license entity.License
	if err := c.ShouldBindJSON(&license); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.licenseService.CreateLicense(c.Request.Context(), &license)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetAllLicenses godoc
// @Summary Get all licenses
// @Description Get list of all clinic licenses
// @Tags licenses
// @Produce json
// @Success 200 {array} entity.License
// @Router /licenses [get]
func (h *LicenseHandler) GetAllLicenses(c *gin.Context) {
	licenses, err := h.licenseService.GetAllLicenses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, licenses)
}

// GetLicenseByID godoc
// @Summary Get license by ID
// @Description Get license details by ID
// @Tags licenses
// @Produce json
// @Param id path int true "License ID"
// @Success 200 {object} entity.License
// @Failure 404 {object} map[string]string
// @Router /licenses/{id} [get]
func (h *LicenseHandler) GetLicenseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	license, err := h.licenseService.GetLicenseByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	c.JSON(http.StatusOK, license)
}

// UpdateLicense godoc
// @Summary Update license
// @Description Update clinic license (admin only)
// @Tags licenses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "License ID"
// @Param request body entity.License true "License update data"
// @Success 200 {object} entity.License
// @Failure 400 {object} map[string]string
// @Router /licenses/{id} [put]
func (h *LicenseHandler) UpdateLicense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	var license entity.License
	if err := c.ShouldBindJSON(&license); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.licenseService.UpdateLicense(c.Request.Context(), id, &license)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteLicense godoc
// @Summary Delete license
// @Description Delete clinic license by ID (admin only)
// @Tags licenses
// @Security BearerAuth
// @Param id path int true "License ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /licenses/{id} [delete]
func (h *LicenseHandler) DeleteLicense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	if err := h.licenseService.DeleteLicense(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
