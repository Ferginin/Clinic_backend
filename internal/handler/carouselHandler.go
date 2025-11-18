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

type CarouselHandler struct {
	carouselService *service.CarouselService
}

func NewCarouselHandler(cfg *config.Config, db *pgxpool.Pool) *CarouselHandler {
	return &CarouselHandler{
		carouselService: service.NewCarouselService(db),
	}
}

// CreateSlide godoc
// @Summary Create carousel slide
// @Description Create a new carousel slide (admin only)
// @Tags carousel
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body entity.Carousel true "Carousel slide data"
// @Success 201 {object} entity.Carousel
// @Failure 400 {object} map[string]string
// @Router /carousel [post]
func (h *CarouselHandler) CreateSlide(c *gin.Context) {
	var carousel entity.Carousel
	if err := c.ShouldBindJSON(&carousel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.carouselService.CreateSlide(c.Request.Context(), &carousel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetAllSlides godoc
// @Summary Get all carousel slides
// @Description Get list of all carousel slides
// @Tags carousel
// @Produce json
// @Success 200 {array} entity.Carousel
// @Router /carousel [get]
func (h *CarouselHandler) GetAllSlides(c *gin.Context) {
	slides, err := h.carouselService.GetAllSlides(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slides)
}

// GetSlideByID godoc
// @Summary Get carousel slide by ID
// @Description Get carousel slide details by ID
// @Tags carousel
// @Produce json
// @Param id path int true "Slide ID"
// @Success 200 {object} entity.Carousel
// @Failure 404 {object} map[string]string
// @Router /carousel/{id} [get]
func (h *CarouselHandler) GetSlideByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slide ID"})
		return
	}

	slide, err := h.carouselService.GetSlideByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slide not found"})
		return
	}

	c.JSON(http.StatusOK, slide)
}

// UpdateSlide godoc
// @Summary Update carousel slide
// @Description Update carousel slide (admin only)
// @Tags carousel
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Slide ID"
// @Param request body entity.Carousel true "Slide update data"
// @Success 200 {object} entity.Carousel
// @Failure 400 {object} map[string]string
// @Router /carousel/{id} [put]
func (h *CarouselHandler) UpdateSlide(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slide ID"})
		return
	}

	var carousel entity.Carousel
	if err := c.ShouldBindJSON(&carousel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.carouselService.UpdateSlide(c.Request.Context(), id, &carousel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteSlide godoc
// @Summary Delete carousel slide
// @Description Delete carousel slide by ID (admin only)
// @Tags carousel
// @Security BearerAuth
// @Param id path int true "Slide ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /carousel/{id} [delete]
func (h *CarouselHandler) DeleteSlide(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slide ID"})
		return
	}

	if err := h.carouselService.DeleteSlide(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
