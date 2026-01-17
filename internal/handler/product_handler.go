package handler

import (
	"net/http"
	"strconv"

	"selfcart/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// GetAll godoc
// @Summary Get all products
// @Description Get all products details
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} []repository.Product
// @Failure 404 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetByBarcode godoc
// @Summary Get product by barcode
// @Description Get product details by its barcode
// @Tags products
// @Accept json
// @Produce json
// @Param barcode path string true "Product Barcode"
// @Success 200 {object} repository.Product
// @Failure 404 {object} map[string]string
// @Router /products/barcode/{barcode} [get]
func (h *ProductHandler) GetByBarcode(c *gin.Context) {
	barcode := c.Param("barcode")

	product, err := h.service.GetByBarcode(c.Request.Context(), barcode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetByCategory godoc
// @Summary Get products by category
// @Description Get products details by its category
// @Tags products
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} []repository.Product
// @Failure 404 {object} map[string]string
// @Router /products/category/{category_id} [get]
func (h *ProductHandler) GetByCategory(c *gin.Context) {
	categoryIDStr := c.Param("category_id")

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
		return
	}

	products, err := h.service.GetByCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
