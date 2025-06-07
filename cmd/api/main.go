package main

import (
	"fmt"
	"github.com/adty404/kredit-plus/internal/platform/database"
	"github.com/adty404/kredit-plus/internal/platform/migration"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, will use OS environment variables")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err := migration.Run(db); err != nil {
		log.Fatalf("Could not run database migrations: %v", err)
	}

	fmt.Println("Application startup completed successfully!")

	router := gin.Default()

	router.GET(
		"/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		},
	)

	log.Println("Starting the HTTP server on port", os.Getenv("SERVE_PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVE_PORT"), router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
