package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"selfcart/internal/service"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *service.CartService
}
type CreateCartRequest struct {
	CustomerID int64 `json:"customer_id" binding:"required"`
}

func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{service: s}
}

// CreateCart godoc
// @Summary Create a new cart
// @Description Create a new cart for a customer
// @Tags cart
// @Accept json
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Success 200 {object} repository.Cart
// @Failure 404 {object} map[string]string
// @Router /cart [post]
func (h *CartHandler) CreateCart(c *gin.Context) {
	var req CreateCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.service.CreateCart(c.Request.Context(), req.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddItem godoc
// @Summary Add an item to the cart
// @Description Add an item to the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param cart_id path int true "Cart ID"
// @Param barcode path string true "Barcode"
// @Param action path string true "Action"
// @Success 200 {object} repository.Cart
// @Failure 404 {object} map[string]string
// @Router /cart/add [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	var req struct {
		CartID  int64  `json:"cart_id" binding:"required"`
		Barcode string `json:"barcode" binding:"required"`
		Action  string `json:"action" binding:"required"`
	}

	fmt.Println(req)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddItem(c.Request.Context(), req.CartID, req.Barcode, req.Action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added to cart"})
}

// RemoveItem godoc
// @Summary Remove an item from the cart
// @Description Remove an item from the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param cart_id path int true "Cart ID"
// @Param item_id path int true "Item ID"
// @Success 200 {object} repository.Cart
// @Failure 404 {object} map[string]string
// @Router /cart/remove [post]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	var req struct {
		CartID int64 `json:"cart_id" binding:"required"`
		ItemID int64 `json:"item_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RemoveItem(c.Request.Context(), req.CartID, req.ItemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart"})
}

// GetCart godoc
// @Summary Get cart details
// @Description Get cart details
// @Tags cart
// @Accept json
// @Produce json
// @Param cart_id path int true "Cart ID"
// @Success 200 {object} repository.Cart
// @Failure 404 {object} map[string]string
// @Router /cart/{cart_id} [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	cartID, err := strconv.ParseInt(c.Param("cart_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart_id"})
		return
	}
	cart, err := h.service.GetCart(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items, err := h.service.GetCartItems(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cart.Items = &items
	cart.Total = cart.Total - cart.Discount
	c.JSON(http.StatusOK, cart)
}

// CheckOut godoc
// @Summary Check out the cart
// @Description Check out the cart
// @Tags cart
// @Accept json
// @Produce json
// @Param cart_id path int true "Cart ID"
// @Success 200 {object} repository.Transaction
// @Failure 404 {object} map[string]string
// @Router /cart/checkout [post]
func (h *CartHandler) CheckOut(c *gin.Context) {
	var req struct {
		CartID int64 `json:"cart_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.service.CheckOut(c.Request.Context(), req.CartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
