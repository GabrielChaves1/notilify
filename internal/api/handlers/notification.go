package handlers

import (
	"GabrielChaves1/notilify/internal/api/middleware"
	"GabrielChaves1/notilify/internal/application/dto/request"
	"GabrielChaves1/notilify/internal/application/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandlers struct {
	createNotificationUseCase *usecase.CreateNotification
}

func NewNotificationHandlers(
	createNotificationUseCase *usecase.CreateNotification,
) *NotificationHandlers {
	return &NotificationHandlers{
		createNotificationUseCase: createNotificationUseCase,
	}
}

func (h *NotificationHandlers) CreateNotification(c *gin.Context) {
	var dto request.CreateNotificationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		middleware.HandleError(c, err)
		return
	}

	result, err := h.createNotificationUseCase.Execute(c.Request.Context(), dto)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Notification created",
		"data":    result,
	})
}
