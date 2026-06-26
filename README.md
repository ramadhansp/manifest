# Manifest & Customs Declaration Service

Backend REST API untuk mengelola proses Container → Manifest → BC 1.1 → NPE → READY XRAY.

## Fitur & Arsitektur
- **Golang 1.24+**
- **Gin Framework** (HTTP Router)
- **GORM & PostgreSQL** (Database ORM)
- **Clean Architecture** (Layer Controller, Service, Repository)
- **Dependency Injection** & **Interface**
- **Context Handling** (Timeout & Context Cancellation)
- **Transaction** (Pembuatan NPE mengubah status Manifest dalam satu transaction)
- **Graceful Shutdown**
- **Middleware Logger & Recovery**

## Struktur Folder
```text
.
├── config/             # Konfigurasi Database (GORM, PostgreSQL)
├── controller/         # Request Binding, Response Formatting
├── dto/                # Data Transfer Object (Request & Response)
├── middleware/         # Gin Middlewares (Logger, Recovery)
├── models/             # Entitas Database (GORM Models)
├── repository/         # Akses Database (Interfaces & Implementasi GORM)
├── service/            # Business Logic & Validation (Beserta Unit Tests)
├── docker-compose.yml  # Docker Compose config
├── Dockerfile          # Multi-stage build Dockerfile
├── main.go             # Entrypoint aplikasi
└── README.md           # Dokumentasi ini
```

## Persyaratan
- Docker & Docker Compose
- Go 1.24+ (Jika ingin menjalankan lokal tanpa Docker)

## Cara Menjalankan

### Menggunakan Docker Compose (Direkomendasikan)
1. Buka terminal di folder project
2. Jalankan perintah:
   ```bash
   docker-compose up --build -d
   ```
3. API akan berjalan di `http://localhost:8080`. PostgreSQL akan otomatis di-setup dan di-migrate (AutoMigrate GORM).

### Menjalankan Lokal (Tanpa Docker)
1. Setup PostgreSQL lokal. Buat database `manifest_db` dengan user `postgres` dan password `postgres`.
2. Install dependensi Golang:
   ```bash
   go mod tidy
   ```
3. Jalankan aplikasi:
   ```bash
   go run main.go
   ```

## Menjalankan Unit Test
```bash
go test ./service -v
```

## Daftar Endpoint & Contoh

Semua request menggunakan standard JSON Response API:
```json
{
  "success": true,
  "message": "Success",
  "data": { ... }
}
```

### 1. Shipping Agents
- `POST /api/shipping-agents`
- `GET /api/shipping-agents`
- `GET /api/shipping-agents/:id`

### 2. Vessels
- `POST /api/vessels`
- `GET /api/vessels`
- `GET /api/vessels/:id`

### 3. Manifests
- `POST /api/manifests`
- `GET /api/manifests?page=1&limit=10&search=MAN`
- `GET /api/manifests/:id`
- `POST /api/manifests/:id/details`

### 4. BC11
- `POST /api/bc11`
  - Validasi: `Container already has active BC11`
  ```json
  {
      "manifest_id": "uuid-here",
      "bc11_number": "BC-12345"
  }
  ```

### 5. NPE
- `POST /api/npe`
  - Validasi: `Invalid BC11`, `Manifest is not complete`
  - Sukses akan mengubah status manifest menjadi `READY_XRAY`.
  ```json
  {
      "bc11_id": "uuid-here",
      "npe_number": "NPE-98765"
  }
  ```

### 6. Summary & Seed
- `GET /api/summary`: Menampilkan statistik jumlah (Aggregate GORM Count).
- `POST /api/seed`: Memicu seeding (Stubs logic implementasi tambahan).
