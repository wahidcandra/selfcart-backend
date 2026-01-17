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

// GetStore godoc
// @Summary Get store details
// @Description Get store details
// @Tags store
// @Accept json
// @Produce json
// @Success 200 {object} repository.Store
// @Failure 404 {object} map[string]string
// @Router /store [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	data, err := h.service.GetStore(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetRack godoc
// @Summary Get rack details
// @Description Get rack details
// @Tags store
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Success 200 {object} []repository.Rack
// @Failure 404 {object} map[string]string
// @Router /store/rack/{store_id} [get]
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

// GetRackZone godoc
// @Summary Get rack zone details
// @Description Get rack zone details
// @Tags store
// @Accept json
// @Produce json
// @Param store_id path int true "Store ID"
// @Param zone path string true "Zone"
// @Success 200 {object} []repository.Rack
// @Failure 404 {object} map[string]string
// @Router /store/rack/zone/{store_id}/{zone} [get]

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
