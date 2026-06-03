package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/promotionos/analytics-service/internal/application"
)

type AnalyticsHandler struct {
	service *application.AnalyticsService
}

func NewAnalyticsHandler(svc *application.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: svc}
}

func (h *AnalyticsHandler) GetReport(c *gin.Context) {
	campaignID := c.Param("id")
	tenantID := c.Query("tenantId")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId required"})
		return
	}
	metrics, err := h.service.GetReport(campaignID, tenantID)
	if err != nil || metrics == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "metrics not found"})
		return
	}
	// BUG: if campaign is 100% Kroger-funded, this will panic in LiftCalculator
	c.JSON(http.StatusOK, metrics)
}

func (h *AnalyticsHandler) GetBurn(c *gin.Context) {
	// TODO Team 5 Sprint 2
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented — Team 5 Sprint 2"})
}

func (h *AnalyticsHandler) GetLift(c *gin.Context) {
	// TODO Team 5 Sprint 3
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented — Team 5 Sprint 3"})
}
