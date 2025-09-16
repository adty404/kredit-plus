# Aplikasi Backend Kredit Plus

Selamat datang di repositori Aplikasi Kredit Plus. Ini adalah layanan backend yang dibangun menggunakan Go (Golang) untuk mensimulasikan sistem aplikasi kredit sederhana. Proyek ini dibangun dengan prinsip Clean Architecture untuk memastikan kode yang bersih, mudah diuji, dan dapat dipelihara.

Aplikasi ini mencakup fungsionalitas inti seperti manajemen konsumen, penetapan limit kredit, pembuatan transaksi, dan sistem otentikasi yang aman.

## ✨ Fitur Utama

* **Manajemen Pengguna & Otentikasi**:
    * Registrasi dan Login untuk pengguna sistem (admin/konsumen).
    * Penggunaan **JWT (JSON Web Tokens)** untuk mengamankan endpoint API.
    * Implementasi keamanan password dengan **hashing bcrypt** (OWASP A07).

* **Manajemen Konsumen**:
    * CRUD (Create, Read, Update, Delete) penuh untuk data konsumen.
    * Upload file untuk foto KTP dan foto selfie saat pendaftaran konsumen.

* **Manajemen Limit Kredit**:
    * Penetapan plafon kredit keseluruhan (`overall_credit_limit`) untuk setiap konsumen.
    * Penetapan batas kredit (`credit_limit`) yang spesifik untuk setiap tenor yang tersedia (1, 2, 3, dan 6 bulan).

* **Manajemen Transaksi**:
    * Pembuatan transaksi kredit dengan validasi terhadap limit tenor dan sisa plafon keseluruhan.
    * Penanganan *race condition* pada saat pembuatan transaksi menggunakan **transaksi database dan pessimistic locking**.

* **Keamanan OWASP Top 10**:
    * ✅ **A01: Broken Access Control**: Rute-rute API diproteksi dengan middleware, memastikan pengguna hanya bisa mengakses data miliknya sendiri.
    * ✅ **A03: Injection**: Aman dari SQL Injection berkat penggunaan GORM dengan *parameterized queries*.
    * ✅ **A07: Identification and Authentication Failures**: Menggunakan hashing bcrypt untuk password dan JWT untuk manajemen sesi.

* **Dockerize**: Seluruh aplikasi dan database-nya telah di-containerize menggunakan Docker dan Docker Compose untuk kemudahan setup dan konsistensi lingkungan.

## 🏗️ Arsitektur & Teknologi

Proyek ini mengadopsi prinsip **Clean Architecture** dengan lapisan-lapisan sebagai berikut:
1.  **Domain**: Berisi model entitas dan kontrak (interface) untuk repository.
2.  **Usecase**: Berisi logika bisnis inti aplikasi.
3.  **Repository**: Implementasi konkret untuk berinteraksi dengan database.
4.  **Handler**: Bertanggung jawab untuk menangani permintaan HTTP, validasi, dan respons.

**Teknologi yang Digunakan:**
* **Bahasa**: Go (Golang) 1.23
* **Web Framework**: Gin
* **Database**: PostgreSQL
* **ORM**: GORM
* **Migrasi**: `golang-migrate/migrate`
* **Validasi**: `go-playground/validator`
* **Containerization**: Docker & Docker Compose

## 📖 Dokumentasi API

Dokumentasi lengkap untuk semua endpoint API, termasuk contoh request dan response, tersedia di Postman. Anda bisa mengaksesnya melalui link di bawah ini.

* **Link Dokumentasi**: [https://documenter.getpostman.com/view/10117768/2sB2x3nDD5](https://documenter.getpostman.com/view/10117768/2sB2x3nDD5)


## 🚀 Cara Menjalankan dengan Docker

Menjalankan proyek ini sangat mudah karena sudah sepenuhnya di-containerize.

### Prasyarat
1.  **Git**: Untuk mengambil kode dari repositori.
2.  **Docker & Docker Compose**: Terinstal di sistem Anda.

### Langkah-langkah

1.  **Clone Repositori**
    ```bash
    git clone <URL_repositori_anda>
    cd <nama_folder_proyek>
    ```

2.  **Buat File `.env`**
    Buat file bernama `.env` di direktori root proyek dan isi dengan konfigurasi berikut. Sesuaikan jika perlu.
    ```env
    # Konfigurasi Database PostgreSQL
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=root
    DB_NAME=kredit_plus
    DB_PORT=5432
    DB_SSLMODE=disable
    DB_TIMEZONE=Asia/Jakarta

    # Konfigurasi Server
    SERVE_PORT=8080

    # Konfigurasi JWT
    JWT_SECRET=kunci_rahasia_yang_sangat_aman
    ```

3.  **Build dan Jalankan Container**
    Jalankan perintah berikut. Perintah ini akan membangun image Docker untuk aplikasi dan database, lalu memulainya.
    ```bash
    docker-compose up --build
    ```
    Biarkan terminal ini berjalan. Aplikasi Anda sekarang bisa diakses di `http://localhost:8080`.

4.  **Jalankan Migrasi & Seeder (di Terminal Baru)**
    Buka **terminal baru**, masuk ke direktori proyek, dan jalankan perintah berikut:

    * **Untuk membuat tabel-tabel di database:**
        ```bash
        docker-compose exec app migrate -database "postgres://postgres:root@db:5432/kredit_plus?sslmode=disable" -path migrations up
        ```
    * **Untuk mengisi data awal (admin, konsumen Budi & Annisa):**
        ```bash
        docker-compose exec app ./kredit-app --seed
        ```

Aplikasi Anda sekarang sudah sepenuhnya siap digunakan!

## ⚙️ Perintah yang Berguna

Semua perintah ini dijalankan dari terminal di root folder proyek Anda.

| Perintah | Deskripsi |
| :--- | :--- |
| `docker-compose up` | Memulai semua layanan (jika image sudah di-build). |
| `docker-compose up -d` | Memulai semua layanan di latar belakang (detached mode). |
| `docker-compose down` | Menghentikan dan menghapus semua container. |
| `docker-compose logs -f app` | Melihat log real-time dari aplikasi Go Anda. |
| `docker-compose exec app make migrate-up` | Menjalankan migrasi UP di dalam container. |
| `docker-compose exec app make migrate-down` | Menjalankan migrasi DOWN di dalam container. |

## 📖 Endpoint API Utama

Berikut adalah daftar endpoint API yang tersedia:

### Otentikasi
* `POST /api/v1/auth/register`
* `POST /api/v1/auth/login`

### Konsumen
* `POST /api/v1/consumers` (Memerlukan otorisasi admin)
* `GET /api/v1/consumers` (Memerlukan otorisasi admin)
* `GET /api/v1/consumers/:id` (Memerlukan autentikasi)
* `PUT /api/v1/consumers/:id` (Memerlukan autentikasi)
* `DELETE /api/v1/consumers/:id` (Memerlukan otorisasi admin)

### Limit Kredit
* `POST /api/v1/consumers/:id/limits` (Memerlukan otorisasi admin)

### Transaksi
* `POST /api/v1/consumers/:id/transactions` (Memerlukan autentikasi)
* `GET /api/v1/consumers/:id/transactions` (Memerlukan autentikasi)
