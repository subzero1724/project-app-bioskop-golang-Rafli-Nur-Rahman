# 🎬 Cinema Booking System API

<div align="center">

**Sistem Pemesanan Bioskop Berbasis RESTful API**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-316192?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

*Aplikasi backend yang dibangun dengan Golang untuk menangani sistem pemesanan tiket bioskop end-to-end*

[Fitur](#-fitur-utama) • [Instalasi](#️-instalasi) • [API Docs](#-dokumentasi-api) • [Arsitektur](#-arsitektur)

</div>

---

## 👤 Info Proyek

- **Nama**: Rafli Nur Rahman
- **Kelas**: Back End Golang
- **Deskripsi**: RESTful API untuk sistem pemesanan tiket bioskop dengan fitur lengkap dari registrasi hingga pembayaran

---

## ✨ Fitur Utama

<table>
<tr>
<td width="50%">

### 🔐 Autentikasi & Keamanan
- ✅ Registrasi pengguna baru
- ✅ Login dengan JWT Token
- ✅ Logout aman
- ✅ Password hashing dengan Bcrypt

</td>
<td width="50%">

### 🎭 Manajemen Bioskop
- ✅ Daftar bioskop tersedia
- ✅ Detail informasi bioskop
- ✅ Pengecekan kursi real-time
- ✅ Sistem booking efisien

</td>
</tr>
<tr>
<td width="50%">

### 💳 Pembayaran
- ✅ Multiple payment methods
- ✅ Integrasi gateway payment
- ✅ Konfirmasi pembayaran
- ✅ Invoice otomatis

</td>
<td width="50%">

### 📊 Riwayat & Laporan
- ✅ History pemesanan
- ✅ Detail transaksi
- ✅ Status booking
- ✅ Export data

</td>
</tr>
</table>

---

## 🛠️ Tech Stack

| Komponen | Teknologi | Versi |
|----------|-----------|-------|
| **Bahasa** | Go | 1.21+ |
| **Web Framework** | Chi Router | Latest |
| **Database** | PostgreSQL | 13+ |
| **Driver DB** | pgx | v5 |
| **Validasi** | go-playground/validator | v10 |
| **Config** | Viper | Latest |
| **Logging** | Uber Zap | Latest |
| **Auth** | JWT (golang-jwt) | v5 |
| **Security** | Bcrypt | Latest |

---

## 📁 Struktur Proyek

```
sistem-pemesanan-bioskop/
│
├── 📂 cmd/
│   └── main.go                    # Entry point aplikasi
│
├── 📂 internal/
│   ├── config/                    # Konfigurasi & environment
│   ├── database/                  # Database connection
│   ├── dto/                       # Data Transfer Objects
│   ├── handler/                   # HTTP Handlers (Controllers)
│   ├── middleware/                # Auth, Logger, CORS middleware
│   ├── models/                    # Domain models & entities
│   ├── repository/                # Data access layer
│   ├── router/                    # API route definitions
│   ├── service/                   # Business logic layer
│   └── utils/                     # Helper functions
│
├── 📂 pkg/
│   └── logger/                    # Custom Zap logger setup
│
├── 📂 migrations/
│   └── 001_init_schema.sql        # Database schema
│
├── .env.example                   # Environment template
├── go.mod                         # Go modules
├── go.sum                         # Dependencies checksum
└── README.md                      # Documentation
```

---

## ⚙️ Instalasi

### 📋 Prasyarat

Pastikan Anda telah menginstall:

- [Go](https://go.dev/dl/) (versi 1.21 atau lebih tinggi)
- [PostgreSQL](https://www.postgresql.org/download/) (versi 13 atau lebih tinggi)
- [Git](https://git-scm.com/downloads)
- API Testing Tool (Postman/Insomnia)

### 🚀 Quick Start

**1. Clone Repository**

```bash
git clone https://github.com/subzero1724/project-app-bioskop-golang-Rafli-Nur-Rahman.git

cd cinema-booking-system
```

**2. Install Dependencies**

```bash
go mod download
go mod tidy
```

**3. Setup Environment**

```bash
cp .env.example .env
# Edit .env sesuai konfigurasi Anda
```

**4. Setup Database**

```bash
# Buat database baru
createdb cinema_booking

# Atau via psql
psql -U postgres
CREATE DATABASE cinema_booking;
\q

# Jalankan migrasi
psql -U postgres -d cinema_booking -f migrations/001_init_schema.sql
```

**5. Run Application**

```bash
# Development mode
go run cmd/main.go

# Build & Run
go build -o bin/cinema-api cmd/main.go
./bin/cinema-api
```

Server akan berjalan di `http://localhost:8080` (default)
git clone <repository-url>
---

## 📚 Dokumentasi API

### Base URL
```
http://localhost:8080/api
```

### 🔑 Authentication Endpoints

<details>
<summary><b>POST</b> <code>/register</code> - Registrasi User Baru</summary>

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "full_name": "John Doe"
  }
}
```
</details>

<details>
<summary><b>POST</b> <code>/login</code> - Login User</summary>

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2026-01-15T12:00:00Z",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com"
    }
  }
}
```

**Headers untuk request selanjutnya:**
```
Authorization: Bearer {token}
```
</details>

<details>
<summary><b>POST</b> <code>/logout</code> - Logout User</summary>

**Headers:**
```
Authorization: Bearer {token}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Logout successful"
}
```
</details>

---

### 🎬 Cinema Endpoints

<details>
<summary><b>GET</b> <code>/cinemas</code> - Daftar Semua Bioskop</summary>

**Query Parameters:**
- `page` (optional): Nomor halaman (default: 1)
- `page_size` (optional): Jumlah item per halaman (default: 10)

**Success Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Cinema XXI Grand Indonesia",
      "location": "Jakarta Pusat",
      "description": "Bioskop premium dengan IMAX dan Dolby Atmos",
      "total_seats": 150,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "Cinema CGV Pacific Place",
      "location": "Jakarta Selatan",
      "description": "Bioskop dengan teknologi 4DX",
      "total_seats": 200,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 10,
    "total_items": 4,
    "total_pages": 1
  }
}
```
</details>

<details>
<summary><b>GET</b> <code>/cinemas/{id}</code> - Detail Bioskop</summary>

**Path Parameters:**
- `id`: ID bioskop

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Cinema XXI Grand Indonesia",
    "location": "Jakarta Pusat",
    "description": "Bioskop premium dengan IMAX dan Dolby Atmos",
    "total_seats": 150,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```
</details>

<details>
<summary><b>GET</b> <code>/cinemas/{id}/seats</code> - Cek Ketersediaan Kursi</summary>

**Path Parameters:**
- `id`: ID bioskop

**Query Parameters:**
- `date` (optional): Tanggal tayang (format: YYYY-MM-DD)
- `time` (optional): Waktu tayang (format: HH:MM)

**Success Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 101,
      "cinema_id": 1,
      "seat_number": "A1",
      "row_number": "A",
      "seat_type": "VIP",
      "price": 75000,
      "is_available": true
    },
    {
      "id": 102,
      "cinema_id": 1,
      "seat_number": "A2",
      "row_number": "A",
      "seat_type": "VIP",
      "price": 75000,
      "is_available": false
    },
    {
      "id": 103,
      "cinema_id": 1,
      "seat_number": "A3",
      "row_number": "A",
      "seat_type": "Regular",
      "price": 50000,
      "is_available": true
    }
  ]
}
```
</details>

---

### 💳 Payment Endpoints

<details>
<summary><b>GET</b> <code>/payments/methods</code> - Daftar Metode Pembayaran</summary>

**Success Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "code": "VA_BCA",
      "name": "BCA Virtual Account",
      "type": "bank_transfer"
    },
    {
      "id": 2,
      "code": "GOPAY",
      "name": "GoPay",
      "type": "e_wallet"
    },
    {
      "id": 3,
      "code": "QRIS",
      "name": "QRIS",
      "type": "qr_code"
    }
  ]
}
```
</details>

---

## 🏗️ Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan pemisahan layer yang jelas:

```
┌─────────────────────────────────────────┐
│         Handler Layer (HTTP)            │
│  (Menerima request, validasi input)     │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Service Layer (Business)        │
│  (Logika bisnis, orchestration)         │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│      Repository Layer (Data Access)     │
│  (Query database, CRUD operations)      │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│           Database (PostgreSQL)         │
└─────────────────────────────────────────┘
```

**Keuntungan:**
- ✅ Testable - Mudah untuk unit testing
- ✅ Maintainable - Kode terorganisir dengan baik
- ✅ Scalable - Mudah dikembangkan
- ✅ Independent - Layer tidak saling bergantung

---

## 🔒 Keamanan

- **JWT Authentication**: Setiap endpoint sensitif dilindungi dengan JWT token
- **Password Hashing**: Menggunakan Bcrypt dengan salt rounds yang aman
- **Input Validation**: Validasi ketat pada setiap request menggunakan validator
- **SQL Injection Prevention**: Menggunakan prepared statements
- **CORS Configuration**: Konfigurasi CORS yang tepat

---

## 📝 Environment Variables

Buat file `.env` berdasarkan `.env.example`:

```env
# Application Configuration
APP_NAME=Cinema Booking System
APP_PORT=8080
APP_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=cinema_booking
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRATION_HOURS=24

# Logging
LOG_LEVEL=info
LOG_ENCODING=json
```

---
## Video Demo:
```
git clone <repository-url>
```
