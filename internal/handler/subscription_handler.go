package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"subscription-service/internal/model"
	"subscription-service/internal/service"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(s *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: s}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка валидации: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.Create(&req)
	if err != nil {
		log.Printf("Ошибка создания: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Создана подписка: ID=%d, Service=%s, User=%s", sub.ID, sub.ServiceName, sub.UserID)
	c.JSON(http.StatusCreated, sub)
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var subscriptionID int
	_, err := fmt.Sscanf(id, "%d", &subscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	sub, err := h.service.GetByID(subscriptionID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "подписка не найдена"})
			return
		}
		log.Printf("Ошибка получения: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var subscriptionID int
	_, err := fmt.Sscanf(id, "%d", &subscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.Update(subscriptionID, &req)
	if err != nil {
		log.Printf("Ошибка обновления: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var subscriptionID int
	_, err := fmt.Sscanf(id, "%d", &subscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID"})
		return
	}

	if err := h.service.Delete(subscriptionID); err != nil {
		log.Printf("Ошибка удаления: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка удаления"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	filters := make(map[string]interface{})

	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if serviceName := c.Query("service_name"); serviceName != "" {
		filters["service_name"] = serviceName
	}
	if startAfter := c.Query("start_after"); startAfter != "" {
		if t, err := time.Parse("01-2006", startAfter); err == nil {
			filters["start_after"] = t
		}
	}
	if endBefore := c.Query("end_before"); endBefore != "" {
		if t, err := time.Parse("01-2006", endBefore); err == nil {
			filters["end_before"] = t
		}
	}

	subs, err := h.service.List(filters)
	if err != nil {
		log.Printf("Ошибка получения списка: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения списка"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	periodStart := c.Query("period_start") // MM-YYYY
	periodEnd := c.Query("period_end")     // MM-YYYY

	if userID == "" || serviceName == "" || periodStart == "" || periodEnd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "требуются user_id, service_name, period_start, period_end"})
		return
	}

	total, err := h.service.CalculateTotalCost(userID, serviceName, periodStart, periodEnd)
	if err != nil {
		log.Printf("Ошибка расчета: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"service_name": serviceName,
		"period_start": periodStart,
		"period_end":   periodEnd,
		"total_cost":   total,
	})
}
