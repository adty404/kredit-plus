package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	httphandler "github.com/adty404/kredit-plus/internal/handler/http"
	"github.com/adty404/kredit-plus/internal/platform/database"
	"github.com/adty404/kredit-plus/internal/platform/migration"
	"github.com/adty404/kredit-plus/internal/platform/seeder"

	"github.com/joho/godotenv"
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

	// 6. Setup Router HTTP
	router := httphandler.SetupRouter(db)

	// 8. Tambahkan route untuk mendapatkan informasi tentang aplikasi
	port := os.Getenv("SERVE_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Application startup completed successfully!")
	log.Printf("Starting the HTTP server on http://localhost:%s\n", port)

	// 9. Mulai HTTP server
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
