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

type DoctorHandler struct {
	doctorService *service.DoctorService
}

func NewDoctorHandler(cfg *config.Config, db *pgxpool.Pool) *DoctorHandler {
	return &DoctorHandler{
		doctorService: service.NewDoctorService(db),
	}
}

// CreateDoctor godoc
// @Summary Create doctor
// @Description Create a new doctor (admin only)
// @Tags doctors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.DoctorCreateRequest true "Doctor data"
// @Success 201 {object} entity.Doctor
// @Failure 400 {object} map[string]string
// @Router /doctors [post]
func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req entity.DoctorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := h.doctorService.CreateDoctor(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doctor)
}

// GetAllDoctors godoc
// @Summary Get all doctors
// @Description Get list of all doctors
// @Tags doctors
// @Produce json
// @Success 200 {array} entity.Doctor
// @Router /doctors [get]
func (h *DoctorHandler) GetAllDoctors(c *gin.Context) {
	doctors, err := h.doctorService.GetAllDoctors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

// GetDoctorByID godoc
// @Summary Get doctor by ID
// @Description Get doctor details by ID
// @Tags doctors
// @Produce json
// @Param id path int true "Doctor ID"
// @Success 200 {object} entity.Doctor
// @Failure 404 {object} map[string]string
// @Router /doctors/{id} [get]
func (h *DoctorHandler) GetDoctorByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	doctor, err := h.doctorService.GetDoctorByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

// GetBySpecialization godoc
// @Summary Get doctors by specialization
// @Description Get doctors filtered by specialization ID
// @Tags doctors
// @Produce json
// @Param id path int true "Specialization ID"
// @Success 200 {array} entity.Doctor
// @Router /doctors/specialization/{id} [get]
func (h *DoctorHandler) GetBySpecialization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid specialization ID"})
		return
	}

	doctors, err := h.doctorService.GetDoctorsBySpecialization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctors)
}

// GetDoctorSchedule godoc
// @Summary Get doctor schedule
// @Description Get schedule for a specific doctor
// @Tags doctors
// @Produce json
// @Param id path int true "Doctor ID"
// @Success 200 {object} entity.Schedule
// @Failure 404 {object} map[string]string
// @Router /doctors/{id}/schedule [get]
func (h *DoctorHandler) GetDoctorSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	schedule, err := h.doctorService.GetDoctorSchedule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// UpdateDoctor godoc
// @Summary Update doctor
// @Description Update doctor information (admin only)
// @Tags doctors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Doctor ID"
// @Param request body entity.DoctorUpdateRequest true "Doctor update data"
// @Success 200 {object} entity.Doctor
// @Failure 400 {object} map[string]string
// @Router /doctors/{id} [put]
func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var req entity.DoctorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := h.doctorService.UpdateDoctor(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

// DeleteDoctor godoc
// @Summary Delete doctor
// @Description Delete doctor by ID (admin only)
// @Tags doctors
// @Security BearerAuth
// @Param id path int true "Doctor ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /doctors/{id} [delete]
func (h *DoctorHandler) DeleteDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	if err := h.doctorService.DeleteDoctor(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
