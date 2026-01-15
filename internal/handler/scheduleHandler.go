package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleService service.ScheduleServiceInterface
}

func NewScheduleHandler(scheduleService service.ScheduleServiceInterface) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

// CreateSchedule godoc
// @Summary Create schedule
// @Description Create a new schedule (admin only)
// @Tags schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.Schedule true "Schedule data"
// @Success 201 {object} entity.Schedule
// @Failure 400 {object} map[string]string
// @Router /schedules [post]
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var schedule entity.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.scheduleService.CreateSchedule(c.Request.Context(), &schedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetAllSchedules godoc
// @Summary Get all schedules
// @Description Get list of all schedules (admin only)
// @Tags schedules
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.Schedule
// @Router /schedules [get]
func (h *ScheduleHandler) GetAllSchedules(c *gin.Context) {
	schedules, err := h.scheduleService.GetAllSchedules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// GetScheduleByID godoc
// @Summary Get schedule by ID
// @Description Get schedule details by ID
// @Tags schedules
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} entity.Schedule
// @Failure 404 {object} map[string]string
// @Router /schedules/{id} [get]
func (h *ScheduleHandler) GetScheduleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	schedule, err := h.scheduleService.GetScheduleByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// GetByDay godoc
// @Summary Get schedules by day
// @Description Get schedules filtered by day of week (1-7)
// @Tags schedules
// @Produce json
// @Param day path int true "Day of week (1=Monday, 7=Sunday)"
// @Success 200 {array} entity.Schedule
// @Router /schedules/day/{day} [get]
func (h *ScheduleHandler) GetByDay(c *gin.Context) {
	day, err := strconv.Atoi(c.Param("day"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day"})
		return
	}

	schedules, err := h.scheduleService.GetScheduleByDay(c.Request.Context(), day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// UpdateSchedule godoc
// @Summary Update schedule
// @Description Update schedule (admin only)
// @Tags schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body entity.Schedule true "Schedule update data"
// @Success 200 {object} entity.Schedule
// @Failure 400 {object} map[string]string
// @Router /schedules/{id} [put]
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var schedule entity.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.scheduleService.UpdateSchedule(c.Request.Context(), id, &schedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteSchedule godoc
// @Summary Delete schedule
// @Description Delete schedule by ID (admin only)
// @Tags schedules
// @Security BearerAuth
// @Param id path int true "Schedule ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /schedules/{id} [delete]
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	if err := h.scheduleService.DeleteSchedule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
