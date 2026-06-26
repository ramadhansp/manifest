# Manifest & Customs Declaration Service API

Backend REST API untuk mengelola proses ekspor-impor mulai dari **Container → Manifest → BC 1.1 → NPE → READY XRAY**. Aplikasi ini dibangun menggunakan prinsip **Clean Architecture**.

---

## 🚀 Teknologi yang Digunakan
- **Golang 1.24+**
- **Gin Framework** (HTTP Router)
- **GORM & PostgreSQL** (Database ORM)
- **JWT (JSON Web Token)** & **Bcrypt** (Autentikasi & Keamanan)
- **Docker & Docker Compose**

---

## 🏗 Struktur Arsitektur (Clean Architecture)
```text
.
├── config/             # Konfigurasi Database & GORM AutoMigrate
├── controller/         # Request Binding & Response Formatting
├── dto/                # Data Transfer Object (Format JSON Masuk/Keluar)
├── middleware/         # Gin Middlewares (Logger, Recovery, CORS, JWT Auth)
├── models/             # Entitas Database (GORM)
├── repository/         # Akses Database (Interfaces & Implementasi)
├── service/            # Business Logic & Validasi Aturan
├── utils/              # Fungsi Pembantu (Generate/Validasi JWT)
├── docker-compose.yml  # Pengaturan Container (API & Postgres DB)
├── Dockerfile          # Multi-stage build Dockerfile (golang:alpine)
└── main.go             # Entrypoint aplikasi (Routing & Auto-Seed)
```

---

## 📖 Manual Book & Business Logic (Aturan Sistem)

Berikut adalah aturan bisnis yang tertanam di dalam Service Layer aplikasi ini:

### 1. Autentikasi dan Keamanan (JWT)
- **Tingkat Akses (Roles)**: Aplikasi mendukung *Role-Based Access Control* dengan hak akses `Administrator`, `Petugas`, atau `Supervisor`.
- **Proses Login**:
  1. Pengguna memanggil endpoint `POST /api/login`.
  2. Sistem memverifikasi *hash* password menggunakan algoritma **Bcrypt**.
  3. Jika sukses, sistem memberikan JWT Token dengan masa aktif 24 Jam.
- **Proteksi API**: Seluruh endpoint (kecuali `/login`) dikunci. Setiap request harus menyertakan Header: `Authorization: Bearer <token>`.

### 2. Manajemen Master Data
- **Shipping Agent (Agen Pelayaran)**: Entitas penanggung jawab kontainer.
- **Vessel (Kapal)**: Memiliki validasi waktu mutlak di mana **Tanggal Keberangkatan (`departure_date`) TIDAK BOLEH lebih awal dari Tanggal Kedatangan (`arrival_date`)**.

### 3. Alur Manifest (Pernyataan Kargo)
- **Status Awal**: Saat Header Manifest dibuat, statusnya adalah **`DRAFT`**.
- **Penambahan Kontainer**: Ketika sebuah detail kontainer ditambahkan ke manifest yang masih `DRAFT`, status manifest secara otomatis akan berubah menjadi **`COMPLETED`**.

### 4. Dokumen Pabean: BC 1.1
- **Validasi Kritis Container**: Sebelum BC 1.1 diterbitkan, sistem akan mengecek semua kontainer di dalam Manifest. **Satu kontainer HANYA BOLEH memiliki maksimal SATU dokumen BC 1.1 yang aktif (`is_active = true`)**. Jika ganda, request ditolak.

### 5. Dokumen Final: NPE (Nota Pelayanan Ekspor)
NPE adalah gerbang akhir persetujuan.
- **Persyaratan**: Harus merujuk pada `BC11ID` yang valid & aktif, serta Manifest yang sudah `COMPLETED` (punya kapal, agen, dan kontainer).
- **ACID Database Transaction**:
  Penerbitan NPE dan perubahan status Manifest dilakukan di dalam satu transaksi database. 
  Saat NPE sukses dibuat, **status Manifest rujukan akan otomatis berubah menjadi `READY_XRAY`**. Jika ada error di tengah jalan, seluruh proses dibatalkan otomatis (*Rollback*).

---

## ⚙️ Cara Menjalankan Aplikasi

### 1. Menggunakan Docker (Sangat Direkomendasikan)
Cukup jalankan satu perintah berikut di terminal:
```bash
docker-compose up --build -d
```
API akan otomatis menyala di `http://localhost:8080` dan Database PostgreSQL menyala di port `5432` internal / port komputer yang Anda set.

*(Catatan: Jika ada error port 8080/5432 bentrok dengan aplikasi lain, ubah angka port sebelah kiri di `docker-compose.yml` Anda).*

### 2. Testing API dengan Postman
File Postman Collection sudah disediakan di dalam folder project:
`Manifest_Customs_API.postman_collection.json`

1. **Import** file tersebut ke aplikasi Postman.
2. Panggil Endpoint **Login (Get Token)**.
3. Postman akan otomatis menyimpan Token JWT Anda. Anda bisa langsung menembak endpoint lainnya tanpa repot *copy-paste* token manual!

---

## 📡 Daftar Endpoint API Utama

1. **Auth & Seed**
   - `POST /api/login` (Body: `username`, `password`)
   - `POST /api/seed` (Memicu pembuatan data statis secara manual)
2. **Master Data**
   - `POST` & `GET /api/shipping-agents`
   - `POST` & `GET /api/vessels`
3. **Manifest & Customs**
   - `POST` & `GET /api/manifests`
   - `POST /api/manifests/:id/details`
   - `POST /api/bc11`
   - `POST /api/npe`
4. **Dashboard**
   - `GET /api/summary` (Menampilkan statistik jumlah dokumen & kontainer)

*(Saat pertama kali `main.go` dijalankan, sistem akan otomatis men-seed user dengan username: **`admin`** dan password: **`admin123`**).*
