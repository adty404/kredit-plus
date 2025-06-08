package http

import (
	"github.com/adty404/kredit-plus/internal/auth"
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
	userRepo := postgres.NewUserRepository(db)

	// Usecase
	consumerUsecase := usecase.NewConsumerUsecase(db, consumerRepo, userRepo)
	consumerCreditLimitUsecase := usecase.NewConsumerCreditLimitUsecase(
		consumerCreditLimitRepo,
		consumerRepo,
	)
	transactionUsecase := usecase.NewTransactionUsecase(
		db,
		transactionRepo,
		consumerRepo,
		consumerCreditLimitRepo,
	)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Handler
	consumerHandler := NewConsumerHandler(consumerUsecase)
	consumerCreditLimitHandler := NewConsumerCreditLimitHandler(
		consumerCreditLimitUsecase,
		consumerUsecase,
	)
	transactionHandler := NewTransactionHandler(transactionUsecase, consumerRepo)
	userHandler := NewUserHandler(userUsecase)

	// === Pendaftaran Rute API ===
	api := router.Group("/api/v1")
	{
		// Grup rute untuk otentikasi (Publik)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
		}

		// Grup rute yang memerlukan autentikasi JWT
		protectedRoutes := api.Group("")
		protectedRoutes.Use(auth.AuthMiddleware())
		{
			// Grup rute untuk consumers di dalam grup terproteksi
			consumerRoutes := protectedRoutes.Group("/consumers")
			{
				// Rute utama untuk consumers
				consumerRoutes.POST("", auth.AuthorizeRole("admin"), consumerHandler.CreateConsumer)
				consumerRoutes.GET("", auth.AuthorizeRole("admin"), consumerHandler.GetAllConsumers)
				consumerRoutes.GET("/:id", consumerHandler.GetConsumerByID)
				consumerRoutes.PUT("/:id", consumerHandler.UpdateConsumer)
				consumerRoutes.DELETE("/:id", auth.AuthorizeRole("admin"), consumerHandler.DeleteConsumer)

				consumerRoutes.POST(
					"/:id/limits",
					auth.AuthorizeRole("admin"),
					consumerCreditLimitHandler.CreateLimitForConsumer,
				)

				consumerRoutes.POST("/:id/transactions", transactionHandler.CreateTransaction)
				consumerRoutes.GET("/:id/transactions", transactionHandler.GetTransactionsByConsumerID)
			}
		}
	}

	return router
}
