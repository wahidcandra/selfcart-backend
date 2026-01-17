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
