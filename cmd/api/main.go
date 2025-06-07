package main

import (
	"flag"
	"fmt"
	"github.com/adty404/kredit-plus/internal/platform/database"
	"github.com/adty404/kredit-plus/internal/platform/migration"
	"github.com/adty404/kredit-plus/internal/platform/seeder"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. Tambahkan flag untuk menjalankan seeder
	runSeeder := flag.Bool("seed", false, "Run the database seeder to populate initial data")
	flag.Parse()

	// 2. Coba memuat file .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, will use OS environment variables")
	}

	// 3. Connect ke database
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// 4. Jalankan migrasi database
	if err := migration.Run(db); err != nil {
		log.Fatalf("Could not run database migrations: %v", err)
	}

	// 5. Cek apakah seeder harus dijalankan
	if *runSeeder {
		seeder.Run(db)
		log.Println("Seeder has been run. Exiting.")
		return
	}

	fmt.Println("Application startup completed successfully!")

	// 6. Inisialisasi router Gin
	router := gin.Default()

	// 7. Tambahkan route untuk health check
	router.GET(
		"/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		},
	)

	// 8. Tambahkan route untuk mendapatkan informasi tentang aplikasi
	port := os.Getenv("SERVE_PORT")
	if port == "" {
		port = "8080"
	}

	// 9. Mulai HTTP server
	log.Println("Starting the HTTP server on port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
