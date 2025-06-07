package http

import (
	"github.com/adty404/kredit-plus/internal/repository/postgres"
	"github.com/adty404/kredit-plus/internal/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET(
		"/ping", func(c *gin.Context) {
			c.JSON(
				http.StatusOK, gin.H{
					"message": "pong",
				},
			)
		},
	)

	// === Dependency Injection ===
	// Repository
	consumerRepo := postgres.NewConsumerRepository(db)

	// Usecase
	consumerUsecase := usecase.NewConsumerUsecase(consumerRepo)

	// Handler
	consumerHandler := NewConsumerHandler(consumerUsecase)

	// === Pendaftaran Rute API ===
	api := router.Group("/api/v1")
	{
		consumerRoutes := api.Group("/consumers")
		{
			consumerRoutes.POST("", consumerHandler.CreateConsumer)
			consumerRoutes.GET("", consumerHandler.GetAllConsumers)
			consumerRoutes.GET("/:id", consumerHandler.GetConsumerByID)
			consumerRoutes.PUT("/:id", consumerHandler.UpdateConsumer)
			consumerRoutes.DELETE("/:id", consumerHandler.DeleteConsumer)
		}
	}

	return router
}
