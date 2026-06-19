package handler

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/service"
	"Clinic_backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CallbackHandler struct {
	callbackService service.CallbackServiceInterface
	auditService    service.AuditServiceInterface
}

func NewCallbackHandler(callbackService service.CallbackServiceInterface, auditService service.AuditServiceInterface) *CallbackHandler {
	return &CallbackHandler{callbackService: callbackService, auditService: auditService}
}

// CreateCallbackRequest godoc
// @Summary Create a callback request
// @Description Submit a request for callback from clinic
// @Tags callback
// @Accept json
// @Produce json
// @Param request body entity.CallbackRequestCreate true "Callback request data"
// @Success 201 {object} entity.CallbackRequest
// @Failure 400 {object} map[string]string
// @Router /callback-requests [post]
func (h *CallbackHandler) CreateCallbackRequest(c *gin.Context) {
	var req entity.CallbackRequestCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	callbackReq, err := h.callbackService.CreateCallbackRequest(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, callbackReq)
}

// GetAllCallbackRequests godoc
// @Summary Get all callback requests (admin only)
// @Description Get list of all callback requests
// @Tags callback
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entity.CallbackRequest
// @Failure 403 {object} map[string]string
// @Router /callback-requests [get]
func (h *CallbackHandler) GetAllCallbackRequests(c *gin.Context) {
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

	requests, total, err := h.callbackService.GetAllCallbackRequestsPaginated(c.Request.Context(), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.PaginatedResponse(c, http.StatusOK, requests, offset, limit, total)
}

// GetCallbackRequestByID godoc
// @Summary Get callback request by ID (admin only)
// @Description Get specific callback request details
// @Tags callback
// @Security BearerAuth
// @Produce json
// @Param id path int true "Callback Request ID"
// @Success 200 {object} entity.CallbackRequest
// @Failure 404 {object} map[string]string
// @Router /callback-requests/{id} [get]
func (h *CallbackHandler) GetCallbackRequestByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid callback request id")
		return
	}

	request, err := h.callbackService.GetCallbackRequestByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, request)
}

// UpdateCallbackRequest godoc
// @Summary Update callback request (admin only)
// @Description Update callback request status
// @Tags callback
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Callback Request ID"
// @Param request body entity.CallbackRequestUpdate true "Update data"
// @Success 200 {object} entity.CallbackRequest
// @Failure 400 {object} map[string]string
// @Router /callback-requests/{id} [put]
func (h *CallbackHandler) UpdateCallbackRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid callback request id")
		return
	}

	var req entity.CallbackRequestUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.callbackService.UpdateCallbackRequest(c.Request.Context(), id, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	adminID := c.GetInt("user_id")
	h.auditService.Log(c.Request.Context(), adminID, "UPDATE", "callback", &id, c.ClientIP())
	utils.SuccessResponse(c, http.StatusOK, updated)
}

// DeleteCallbackRequest godoc
// @Summary Delete callback request (admin only)
// @Description Delete callback request
// @Tags callback
// @Security BearerAuth
// @Param id path int true "Callback Request ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /callback-requests/{id} [delete]
func (h *CallbackHandler) DeleteCallbackRequest(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid callback request id")
		return
	}

	if err := h.callbackService.DeleteCallbackRequest(c.Request.Context(), id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	adminID := c.GetInt("user_id")
	h.auditService.Log(c.Request.Context(), adminID, "DELETE", "callback", &id, c.ClientIP())
	c.Status(http.StatusNoContent)
}
