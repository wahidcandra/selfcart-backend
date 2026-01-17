package handler

import (
	"net/http"

	"selfcart/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// GetAll godoc
// @Summary Get all categories
// @Description Get all categories details
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} []repository.Category
// @Failure 404 {object} map[string]string
// @Router /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
