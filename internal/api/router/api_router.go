package router

import (
	"GabrielChaves1/notilify/internal/api/handlers"
	"GabrielChaves1/notilify/internal/api/middleware"
	"GabrielChaves1/notilify/internal/domain/types"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type APIRouterConfig struct {
	Environment types.Environment
}

func SetupAPIRouter(
	config APIRouterConfig,
	logger *logrus.Entry,
	notificationHandlers *handlers.NotificationHandlers,
) *gin.Engine {
	if config.Environment == types.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(middleware.CurrentTime())
	router.Use(middleware.ErrorHandler())
	router.Use(requestid.New())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	notifications := router.Group("/notifications")
	{
		notifications.POST("", notificationHandlers.CreateNotification)
	}
	return router
}
