# Makefile - Disesuaikan untuk bekerja di berbagai platform (termasuk Windows)

# --- Variabel Konfigurasi ---
# Nama file binary yang akan dihasilkan
BINARY_NAME=kredit-app.exe
# Direktori untuk menyimpan hasil build
BUILD_DIR=bin
# Path ke file main.go
MAIN_FILE=cmd/api/main.go
# Path ke folder migrasi
MIGRATION_PATH=migrations
# Koneksi database dari environment variable, dengan fallback jika tidak ada
DATABASE_URL ?= "postgres://postgres:root@localhost:5432/kredit_plus?sslmode=disable"


# --- Perintah Utama ---

# Target default: menjalankan aplikasi menggunakan 'go run'
run:
	@echo "Running the application..."
	@go run $(MAIN_FILE)

run-seed:
	@echo "Running the application with seed data..."
	@go run $(MAIN_FILE) --seed
# docker-compose exec app ./kredit-app --seed

# Build aplikasi
build:
	@echo "Building the application..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

# Menjalankan binary yang sudah di-build
run-build: build
	@echo "Running the built binary..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Membersihkan hasil build
clean:
	@echo "Cleaning build directory..."
	@if exist $(BUILD_DIR) (rmdir /s /q $(BUILD_DIR))

# Menjalankan tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Memformat kode
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Menjalankan linter
lint:
	@echo "Linting code..."
	@golangci-lint run

# Menginstal dependensi
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# --- Perintah Migrasi Database ---
# Catatan: Perintah ini menggunakan migrate CLI. Pastikan sudah terinstall.

# Menjalankan migrasi UP
migrate-up:
	@echo "Running database migrations up..."
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_PATH) up

# docker-compose exec app migrate -database "postgres://postgres:root@db:5432/kredit_plus?sslmode=disable" -path migrations up

# Menjalankan migrasi DOWN
migrate-down:
	@echo "Rolling back last migration..."
	@migrate -database "$(DATABASE_URL)" -path migrations down 1

# docker-compose exec app migrate -database "postgres://postgres:root@db:5432/kredit_plus?sslmode=disable" -path migrations down 1

# Membuat file migrasi baru
# Cara penggunaan: make migrate-create name=nama_migrasi_anda
migrate-create:
	@echo "Creating new migration files..."
	@if [ -z "$(name)" ]; then echo "Error: 'name' is a required argument. Usage: make migrate-create name=your_migration_name"; exit 1; fi
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)


# --- Bantuan ---
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application using 'go run' (default)"
	@echo "  make build        - Build the application"
	@echo "  make run-build    - Run the built binary"
	@echo "  make clean        - Clean build directory"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format the code"
	@echo "  make lint         - Run linter"
	@echo "  make deps         - Install dependencies"
	@echo "  make migrate-up   - Run database migrations up"
	@echo "  make migrate-down - Rollback the last migration"
	@echo "  make migrate-create name=<name> - Create a new migration file"
	@echo "  make help         - Show this help message"

# Default target jika `make` dijalankan tanpa argumen
.DEFAULT_GOAL := run

# Mengabaikan nama file yang sama dengan target make
.PHONY: run build run-build clean test fmt lint deps migrate-up migrate-down migrate-create help