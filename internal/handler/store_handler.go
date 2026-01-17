package handler

import (
	"net/http"
	"strconv"

	"selfcart/internal/service"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	service *service.StoreService
}

func NewStoreHandler(s *service.StoreService) *StoreHandler {
	return &StoreHandler{service: s}
}

func (h *StoreHandler) GetStore(c *gin.Context) {
	data, err := h.service.GetStore(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *StoreHandler) GetRack(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("store_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.GetRack(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *StoreHandler) GetRackZone(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("store_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zone := c.Param("zone")
	data, err := h.service.GetRackZone(c.Request.Context(), storeID, zone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
