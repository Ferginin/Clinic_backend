package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/service"
	"Clinic_backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorService      service.DoctorServiceInterface
	appointmentService service.AppointmentServiceInterface
	auditService       service.AuditServiceInterface
}

func NewDoctorHandler(doctorService service.DoctorServiceInterface, appointmentService service.AppointmentServiceInterface, auditService service.AuditServiceInterface) *DoctorHandler {
	return &DoctorHandler{
		doctorService:      doctorService,
		appointmentService: appointmentService,
		auditService:       auditService,
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
// @Failure 400 {object} utils.ErrorResponse
// @Router /doctors [post]
func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req entity.DoctorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := utils.ValidateStringLength("fullname", req.Fullname, 2, 200); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	doctor, err := h.doctorService.CreateDoctor(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create doctor", err.Error())
		return
	}

	adminID := c.GetInt("user_id")
	h.auditService.Log(c.Request.Context(), adminID, "CREATE", "doctor", &doctor.ID, c.ClientIP())
	utils.SuccessResponse(c, http.StatusCreated, doctor)
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get doctors", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, doctors)
}

// GetDoctorByID godoc
// @Summary Get doctor by ID
// @Description Get doctor details by ID
// @Tags doctors
// @Produce json
// @Param id path int true "Doctor ID"
// @Success 200 {object} entity.Doctor
// @Failure 404 {object} utils.ErrorResponse
// @Router /doctors/{id} [get]
func (h *DoctorHandler) GetDoctorByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	doctor, err := h.doctorService.GetDoctorByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Doctor not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, doctor)
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid specialization ID")
		return
	}

	doctors, err := h.doctorService.GetDoctorsBySpecialization(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get doctors", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, doctors)
}

// GetDoctorSchedule godoc
// @Summary Get doctor schedule
// @Description Get schedule for a specific doctor
// @Tags doctors
// @Produce json
// @Param id path int true "Doctor ID"
// @Success 200 {object} entity.Schedule
// @Failure 404 {object} utils.ErrorResponse
// @Router /doctors/{id}/schedule [get]
func (h *DoctorHandler) GetDoctorSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	schedule, err := h.doctorService.GetDoctorSchedule(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, schedule)
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
// @Failure 400 {object} utils.ErrorResponse
// @Router /doctors/{id} [put]
func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	var req entity.DoctorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if req.Fullname != nil {
		if err := utils.ValidateStringLength("fullname", *req.Fullname, 2, 200); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	doctor, err := h.doctorService.UpdateDoctor(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update doctor", err.Error())
		return
	}

	adminID := c.GetInt("user_id")
	h.auditService.Log(c.Request.Context(), adminID, "UPDATE", "doctor", &id, c.ClientIP())
	utils.SuccessResponse(c, http.StatusOK, doctor)
}

// DeleteDoctor godoc
// @Summary Delete doctor
// @Description Delete doctor by ID (admin only)
// @Tags doctors
// @Security BearerAuth
// @Param id path int true "Doctor ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Router /doctors/{id} [delete]
func (h *DoctorHandler) DeleteDoctor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID")
		return
	}

	if err := h.doctorService.DeleteDoctor(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete doctor", err.Error())
		return
	}

	adminID := c.GetInt("user_id")
	h.auditService.Log(c.Request.Context(), adminID, "DELETE", "doctor", &id, c.ClientIP())
	c.Status(http.StatusNoContent)
}

// GetMyProfile godoc
// @Summary Get own doctor profile
// @Description Get doctor profile linked to current user (doctor only)
// @Tags doctors
// @Security BearerAuth
// @Produce json
// @Success 200 {object} entity.Doctor
// @Failure 404 {object} utils.ErrorResponse
// @Router /doctors/me [get]
func (h *DoctorHandler) GetMyProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	doctor, err := h.doctorService.GetDoctorByUserID(c.Request.Context(), userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "doctor profile not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, doctor)
}

// GetMyScheduleHandler godoc
// @Summary Get own schedule
// @Description Get schedule for current doctor (doctor only)
// @Tags doctors
// @Security BearerAuth
// @Produce json
// @Success 200 {object} entity.Schedule
// @Failure 404 {object} utils.ErrorResponse
// @Router /doctors/me/schedule [get]
func (h *DoctorHandler) GetMyScheduleHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	schedule, err := h.doctorService.GetMySchedule(c.Request.Context(), userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, schedule)
}

// GetMyAppointments godoc
// @Summary Get own appointments
// @Description Get appointments for current doctor (doctor only)
// @Tags doctors
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.Appointment
// @Failure 404 {object} utils.ErrorResponse
// @Router /doctors/me/appointments [get]
func (h *DoctorHandler) GetMyAppointmentsHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	doctor, err := h.doctorService.GetDoctorByUserID(c.Request.Context(), userID.(int))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "doctor profile not found")
		return
	}

	appointments, err := h.appointmentService.GetAppointmentsByDoctor(c.Request.Context(), doctor.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, appointments)
}
