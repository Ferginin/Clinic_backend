package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/service"
	"net/http"
	"strconv"

	"Clinic_backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentService service.AppointmentServiceInterface
}

func NewAppointmentHandler(appointmentService service.AppointmentServiceInterface) *AppointmentHandler {
	return &AppointmentHandler{appointmentService: appointmentService}
}

// CreateAppointment godoc
// @Summary Create an appointment (patient)
// @Description Schedule a new appointment with doctor
// @Tags appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.AppointmentCreateRequest true "Appointment data"
// @Success 201 {object} entity.Appointment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /appointments [post]
func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user_id")
		return
	}

	var req entity.AppointmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	appointment, err := h.appointmentService.CreateAppointment(c.Request.Context(), userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, appointment)
}

// GetMyAppointments godoc
// @Summary Get own appointments (patient)
// @Tags appointments
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.Appointment
// @Failure 401 {object} map[string]string
// @Router /appointments/me [get]
func (h *AppointmentHandler) GetMyAppointments(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid user_id")
		return
	}

	appointments, err := h.appointmentService.GetAppointmentsByPatient(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, appointments)
}

// GetDoctorAppointments godoc
// @Summary Get appointments for doctor (doctor/admin)
// @Tags appointments
// @Security BearerAuth
// @Produce json
// @Param doctor_id path int true "Doctor ID"
// @Success 200 {array} entity.Appointment
// @Failure 401 {object} map[string]string
// @Router /appointments/doctor/{doctor_id} [get]
func (h *AppointmentHandler) GetDoctorAppointments(c *gin.Context) {
	roleValue, exists := c.Get("role")
	if !exists {
		utils.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}
	role, ok := roleValue.(string)
	if !ok || (role != "doctor" && role != "admin") {
		utils.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}

	doctorID, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid doctor id")
		return
	}

	appointments, err := h.appointmentService.GetAppointmentsByDoctor(c.Request.Context(), doctorID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, appointments)
}

// UpdateAppointment godoc
// @Summary Update appointment (patient/admin)
// @Tags appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Appointment ID"
// @Param request body entity.AppointmentUpdateRequest true "Appointment update data"
// @Success 200 {object} entity.Appointment
// @Failure 400 {object} map[string]string
// @Router /appointments/{id} [put]
func (h *AppointmentHandler) UpdateAppointment(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	userID := userIDValue.(int)

	roleValue, _ := c.Get("role")
	role := roleValue.(string)

	appointmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid appointment id")
		return
	}

	var req entity.AppointmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.appointmentService.UpdateAppointment(c.Request.Context(), appointmentID, userID, role, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, updated)
}

// CancelAppointment godoc
// @Summary Cancel appointment (patient/admin)
// @Tags appointments
// @Security BearerAuth
// @Param id path int true "Appointment ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /appointments/{id} [delete]
func (h *AppointmentHandler) CancelAppointment(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	userID := userIDValue.(int)

	roleValue, _ := c.Get("role")
	role := roleValue.(string)

	appointmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid appointment id")
		return
	}

	if err := h.appointmentService.CancelAppointment(c.Request.Context(), appointmentID, userID, role); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// AddAppointmentResult godoc
// @Summary Doctor adds result for appointment
// @Tags appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Appointment ID"
// @Param request body entity.AppointmentResultRequest true "Result data"
// @Success 200 {object} entity.Appointment
// @Failure 400 {object} map[string]string
// @Router /appointments/{id}/result [post]
func (h *AppointmentHandler) AddAppointmentResult(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	userID := userIDValue.(int)

	roleValue, _ := c.Get("role")
	role := roleValue.(string)
	if role != "doctor" && role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}

	appointmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid appointment id")
		return
	}

	var req entity.AppointmentResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.appointmentService.AddAppointmentResult(c.Request.Context(), appointmentID, userID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, updated)
}

// GetAllAppointments godoc
// @Summary Get all appointments (admin)
// @Tags appointments
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.Appointment
// @Failure 403 {object} map[string]string
// @Router /appointments [get]
func (h *AppointmentHandler) GetAllAppointments(c *gin.Context) {
	roleValue, _ := c.Get("role")
	role := roleValue.(string)
	if role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}

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

	appointments, total, err := h.appointmentService.GetAllAppointmentsPaginated(c.Request.Context(), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.PaginatedResponse(c, http.StatusOK, appointments, offset, limit, total)
}

// AdminCreateAppointment godoc
// @Summary Admin creates appointment for a patient or guest (admin)
// @Description Admin books an appointment on behalf of an existing patient or a guest caller
// @Tags appointments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.AppointmentAdminCreateRequest true "Admin appointment data"
// @Success 201 {object} entity.Appointment
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /appointments/admin [post]
func (h *AppointmentHandler) AdminCreateAppointment(c *gin.Context) {
	roleValue, exists := c.Get("role")
	if !exists {
		utils.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}
	role, ok := roleValue.(string)
	if !ok || role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	adminID := userIDValue.(int)

	var req entity.AppointmentAdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	appointment, err := h.appointmentService.AdminCreateAppointment(c.Request.Context(), adminID, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, appointment)
}

// GetAvailableSlots godoc
// @Summary Get available time slots
// @Description Get available appointment slots for a doctor on a specific date
// @Tags appointments
// @Produce json
// @Param doctor_id path int true "Doctor ID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} entity.AvailableSlotsResponse
// @Failure 400 {object} utils.Response
// @Router /appointments/slots/{doctor_id} [get]
func (h *AppointmentHandler) GetAvailableSlots(c *gin.Context) {
	doctorID, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid doctor_id")
		return
	}

	date := c.Query("date")
	if date == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "date query parameter is required (YYYY-MM-DD)")
		return
	}

	slots, err := h.appointmentService.GetAvailableSlots(c.Request.Context(), doctorID, date)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, slots)
}
