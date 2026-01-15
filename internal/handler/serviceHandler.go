package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceService service.ServiceServiceInterface
}

func NewServiceHandler(serviceService service.ServiceServiceInterface) *ServiceHandler {
	return &ServiceHandler{
		serviceService: serviceService,
	}
}

// CreateService godoc
// @Summary Create service
// @Description Create a new medical service (admin only)
// @Tags services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.ServiceCreateRequest true "Service data"
// @Success 201 {object} entity.Service
// @Failure 400 {object} map[string]string
// @Router /services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var req entity.ServiceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc, err := h.serviceService.CreateService(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, svc)
}

// GetAllServices godoc
// @Summary Get all services
// @Description Get list of all medical services
// @Tags services
// @Produce json
// @Success 200 {array} entity.Service
// @Router /services [get]
func (h *ServiceHandler) GetAllServices(c *gin.Context) {
	services, err := h.serviceService.GetAllServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetServiceByID godoc
// @Summary Get service by ID
// @Description Get service details by ID
// @Tags services
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} entity.Service
// @Failure 404 {object} map[string]string
// @Router /services/{id} [get]
func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	svc, err := h.serviceService.GetServiceByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, svc)
}

// GetByCategory godoc
// @Summary Get services by category
// @Description Get services filtered by category ID
// @Tags services
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {array} entity.Service
// @Router /services/category/{id} [get]
func (h *ServiceHandler) GetByCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	services, err := h.serviceService.GetServicesByCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetBySpecialization godoc
// @Summary Get services by specialization
// @Description Get services filtered by specialization ID
// @Tags services
// @Produce json
// @Param id path int true "Specialization ID"
// @Success 200 {array} entity.Service
// @Router /services/specialization/{id} [get]
func (h *ServiceHandler) GetBySpecialization(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid specialization ID"})
		return
	}

	services, err := h.serviceService.GetServicesBySpecialization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// UpdateService godoc
// @Summary Update service
// @Description Update service information (admin only)
// @Tags services
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param request body entity.ServiceCreateRequest true "Service update data"
// @Success 200 {object} entity.Service
// @Failure 400 {object} map[string]string
// @Router /services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	var req entity.ServiceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	svc, err := h.serviceService.UpdateService(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, svc)
}

// DeleteService godoc
// @Summary Delete service
// @Description Delete service by ID (admin only)
// @Tags services
// @Security BearerAuth
// @Param id path int true "Service ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	if err := h.serviceService.DeleteService(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
