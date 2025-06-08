package http

import (
	"github.com/adty404/kredit-plus/internal/repository/postgres"
	"github.com/adty404/kredit-plus/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strings"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Daftarkan Validator untuk binding
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(
			func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			},
		)
	}

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
	consumerCreditLimitRepo := postgres.NewConsumerCreditLimitRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	// Usecase
	consumerUsecase := usecase.NewConsumerUsecase(consumerRepo)
	consumerCreditLimitUsecase := usecase.NewConsumerCreditLimitUsecase(
		consumerCreditLimitRepo,
		consumerRepo,
	)
	transactionUsecase := usecase.NewTransactionUsecase(
		transactionRepo,
		consumerRepo,
		consumerCreditLimitRepo,
	)

	// Handler
	consumerHandler := NewConsumerHandler(consumerUsecase)
	consumerCreditLimitHandler := NewConsumerCreditLimitHandler(
		consumerCreditLimitUsecase,
		consumerUsecase,
	)
	transactionHandler := NewTransactionHandler(transactionUsecase)

	// === Pendaftaran Rute API ===
	api := router.Group("/api/v1")
	{
		consumerRoutes := api.Group("/consumers")
		{
			// Rute utama untuk consumers
			consumerRoutes.POST("", consumerHandler.CreateConsumer)
			consumerRoutes.GET("", consumerHandler.GetAllConsumers)
			consumerRoutes.GET("/:id", consumerHandler.GetConsumerByID)
			consumerRoutes.PUT("/:id", consumerHandler.UpdateConsumer)
			consumerRoutes.DELETE("/:id", consumerHandler.DeleteConsumer)

			consumerRoutes.POST("/:id/limits", consumerCreditLimitHandler.CreateLimitForConsumer)

			consumerRoutes.POST("/:id/transactions", transactionHandler.CreateTransaction)
			consumerRoutes.GET("/:id/transactions", transactionHandler.GetTransactionsByConsumerID)
		}
	}

	return router
}
