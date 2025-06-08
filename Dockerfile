# --- Tahap 1: Build Stage ---
FROM golang:1.23-alpine AS builder

# Instal dependensi yang dibutuhkan: git untuk go mod, tzdata, dan wget untuk mengunduh migrate
RUN apk add --no-cache git tzdata wget

RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.linux-amd64.tar.gz && \
    mv migrate /usr/local/bin/migrate

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/kredit-app ./cmd/api

# --- Tahap 2: Final Stage ---
FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/out/kredit-app .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY migrations ./migrations

# Expose port yang akan digunakan oleh aplikasi Anda.
EXPOSE 8080

# Command untuk menjalankan aplikasi saat container dimulai.
CMD ["./kredit-app"]