# TPS-PKB — Sistem Pemrosesan Pajak Kendaraan Bermotor

REST API backend untuk pengelolaan Pajak Kendaraan Bermotor (PKB) di lingkungan Samsat.  
Dibangun sebagai take-home test dengan arsitektur Hexagonal (Ports & Adapters).

---

## Daftar Isi

- [Fitur Utama](#fitur-utama)
- [Teknologi](#teknologi)
- [Arsitektur](#arsitektur)
- [Alur Request (End-to-End)](#alur-request-end-to-end)
- [Prasyarat](#prasyarat)
- [Mulai Cepat](#mulai-cepat)
- [Konfigurasi Lingkungan](#konfigurasi-lingkungan)
- [Endpoint API](#endpoint-api)
- [Dokumentasi API (Swagger)](#dokumentasi-api-swagger)
- [Aturan Denda](#aturan-denda)
- [Peran & Akses](#peran--akses)
- [Menjalankan Tes](#menjalankan-tes)
- [Perintah Pengembang](#perintah-pengembang)

---

## Fitur Utama

- Manajemen data **Wajib Pajak** (individu dengan NIK, badan usaha dengan NPWP/NIB)
- Pencatatan **Kendaraan** dan **Kewajiban Pajak** tahunan per kendaraan
- Transaksi **Pembayaran** dengan kalkulasi denda otomatis:
  - **Denda Telat Bayar**: 2% × pokok pajak
  - **Denda Kurang Bayar**: 1% × selisih kekurangan
  - **Denda Gabungan**: keduanya jika terlambat sekaligus kurang bayar
- **Laporan** rekapitulasi per periode (total kewajiban, dibayar, denda, sisa)
- **Audit Trail** append-only untuk setiap mutasi data (CREATE/UPDATE)
- **JWT-based RBAC** dengan tiga peran: ADMIN, PETUGAS, WAJIB_PAJAK
- Migrasi database otomatis saat startup

---

## Teknologi

| Komponen       | Library                                |
|----------------|----------------------------------------|
| Bahasa         | Go 1.22 (stdlib `net/http`)            |
| Database       | PostgreSQL 16 + `pgx/v5` + pgxpool     |
| Autentikasi    | JWT HS256 (`golang-jwt/jwt/v5`)        |
| Password       | bcrypt cost 12 (`golang.org/x/crypto`) |
| Nilai uang     | `shopspring/decimal` (bukan float64)   |
| Validasi       | `go-playground/validator/v10`          |
| Query builder  | `squirrel` (filter dinamis)            |
| Migrasi DB     | `golang-migrate/migrate/v4`            |
| Konfigurasi    | `viper` (file `.env` + env var)        |
| Mock testing   | `go.uber.org/mock/mockgen`             |

---

## Arsitektur

```
cmd/server/main.go              ← hanya wire-up dependensi
internal/domain/
  entity/                       ← struct murni, tanpa impor adapter
  dto/                          ← tipe request/filter/error domain
  port/                         ← interface: service, repository, audit, token
  service/                      ← logika bisnis (hanya bergantung pada port)
internal/adapter/
  http/                         ← handler, middleware, router, response
  postgres/                     ← implementasi repository (pgx)
internal/infrastructure/
  config/                       ← binding Viper
  database/                     ← pgxpool + runner migrasi
  token/                        ← implementasi JWT TokenProvider
internal/mock/                  ← mock yang di-generate (commit ke repo)
migrations/                     ← file .up.sql / .down.sql bernomor
```

**Aturan utama:** `internal/domain/` tidak boleh mengimpor dari `adapter/` atau `infrastructure/`.

---

## Alur Request (End-to-End)

### 1. Alur Umum — Semua Protected Endpoint

Setiap request melewati tiga lapis sebelum menyentuh logika bisnis:
**Middleware → Handler → Service → Repository → Database**

```mermaid
sequenceDiagram
    autonumber
    participant C  as Client
    participant R  as Router<br/>(net/http ServeMux)
    participant AM as Auth Middleware<br/>(JWT Validate)
    participant RM as RequireRole<br/>Middleware
    participant SO as SelfOnly<br/>Middleware *
    participant H  as Handler<br/>(adapter/http)
    participant S  as Service<br/>(domain/service)
    participant Re as Repository<br/>(adapter/postgres)
    participant DB as PostgreSQL

    C->>R: HTTP Request + Bearer Token

    alt Endpoint publik (login, health, swagger)
        R-->>C: Response langsung (tanpa middleware)
    end

    R->>AM: Teruskan ke middleware chain
    AM->>AM: Ekstrak & validasi JWT
    alt Token tidak ada / tidak valid
        AM-->>C: 401 UNAUTHORIZED
    end
    AM->>AM: Inject claims + IP + User-Agent ke context
    AM->>RM: next.ServeHTTP()

    RM->>RM: Cek role dari context claims
    alt Role tidak diizinkan
        RM-->>C: 403 FORBIDDEN
    end
    RM->>SO: next.ServeHTTP()

    Note over SO: * Hanya pada endpoint list<br/>(wajib-pajak & kewajiban-pajak)
    SO->>SO: Jika role = WAJIB_PAJAK,<br/>inject owner_filter UUID ke context
    SO->>H: next.ServeHTTP()

    H->>H: Parse & decode JSON body
    H->>H: Validasi struct (validator/v10)
    alt Validasi gagal
        H-->>C: 400 VALIDATION_ERROR
    end
    H->>S: Panggil method service (via interface port)

    S->>S: Jalankan logika bisnis &<br/>aturan domain
    alt Aturan domain dilanggar
        S-->>H: error domain (ErrWajibPajakInaktif, dll)
        H-->>C: 422 UNPROCESSABLE
    end
    S->>Re: Panggil method repository (via interface port)

    Re->>DB: SQL query via pgx/v5
    DB-->>Re: Result rows
    Re-->>S: Entity / error

    S->>S: Tulis audit_log (async, tidak gagalkan request)
    S-->>H: Entity hasil / error

    H->>H: Bungkus dalam response envelope
    H-->>C: 200/201 JSON Response
```

---

### 2. Alur Khusus — POST /api/v1/pembayaran

Endpoint pembayaran memiliki alur paling kompleks: kalkulasi denda + penyimpanan atomik dalam satu transaksi database.

```mermaid
flowchart TD
    A([Client: POST /api/v1/pembayaran]) --> B[Auth Middleware\nValidasi JWT]
    B -->|token tidak valid| ERR1([401 Unauthorized])
    B --> C[RequireRole\nADMIN atau PETUGAS]
    C -->|role tidak sesuai| ERR2([403 Forbidden])
    C --> D[PembayaranHandler.Create\nParse + Validasi request]
    D -->|validasi gagal| ERR3([400 Validation Error])
    D --> E[PembayaranService.Create]

    E --> F{jumlah_bayar\n> 0?}
    F -->|tidak| ERR4([422 Jumlah Harus Positif])
    F -->|ya| G{tanggal_bayar\n≤ hari ini?}
    G -->|tidak| ERR5([422 Tanggal Masa Depan])
    G -->|ya| H[KewajibanRepository\nFindByID]
    H -->|not found| ERR6([404 Not Found])
    H --> I[WajibPajakRepository\nFindByID]
    I --> J{is_active?}
    J -->|tidak| ERR7([422 Wajib Pajak Inaktif])
    J -->|ya| K[DendaService.Hitung\nPure function — tanpa I/O]

    K --> L{tanggal_bayar\n> periode_final?}
    L -->|ya| M[Buat Denda\nTELAT_BAYAR\n2% × pokok]
    K --> N{total_dibayar + jumlah\n< pokok_pajak?}
    N -->|ya| O[Buat Denda\nKURANG_BAYAR\n1% × selisih]

    M --> P[Hitung Status Kewajiban]
    O --> P
    L -->|tidak| P
    N -->|tidak| P

    P --> Q{total_baru\nvs pokok}
    Q -->|=| R1[Status: LUNAS]
    Q -->|kurang| R2[Status: KURANG_BAYAR]
    Q -->|lebih| R3[Status: LEBIH_BAYAR]

    R1 --> S[PembayaranRepository\nSaveWithDenda]
    R2 --> S
    R3 --> S

    S --> T[(PostgreSQL\nSatu Transaksi DB)]
    T --> T1[INSERT pembayaran]
    T --> T2[INSERT denda\n0–2 rows]
    T --> T3[UPDATE kewajiban_pajak\ntotal_dibayar + status]
    T1 & T2 & T3 --> U{Commit OK?}
    U -->|gagal| ERR8([500 Internal Error\nRollback otomatis])
    U -->|berhasil| V[AuditService.Log\nTulis audit_logs]
    V --> W([201 Created\nPembayaran + Denda])
```

---

### 3. Layer Dependency

Panah menunjukkan arah impor. `domain/` tidak pernah mengimpor ke luar.

```
┌─────────────────────────────────────────────────────┐
│                    CLIENT                           │
└──────────────────────┬──────────────────────────────┘
                       │ HTTP
┌──────────────────────▼──────────────────────────────┐
│          adapter/http  (Driving Adapter)             │
│  Router → Middleware → Handler → response.Write*()  │
└──────────────────────┬──────────────────────────────┘
                       │ memanggil via port interface
┌──────────────────────▼──────────────────────────────┐
│               domain/service                        │
│  Logika bisnis, kalkulasi denda, aturan domain      │
│  Hanya bergantung pada interface di domain/port/    │
└──────────┬───────────────────────────┬──────────────┘
           │ port.Repository           │ port.AuditService
┌──────────▼──────────┐   ┌───────────▼──────────────┐
│  adapter/postgres   │   │   adapter/postgres        │
│  (Driven Adapter)   │   │   AuditRepository         │
│  pgx queries        │   │   INSERT audit_logs       │
└──────────┬──────────┘   └───────────┬───────────────┘
           │                          │
┌──────────▼──────────────────────────▼───────────────┐
│                   PostgreSQL 16                      │
└──────────────────────────────────────────────────────┘
```

---

## Prasyarat

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- Go 1.22+ (hanya untuk pengembangan lokal tanpa Docker)

---

## Mulai Cepat

### Menggunakan Docker Compose (Direkomendasikan)

```bash
# 1. Salin file konfigurasi
cp .env.example .env

# 2. Jalankan seluruh stack (app + PostgreSQL)
docker-compose up --build

# Server berjalan di http://localhost:8080
```

Migrasi database dan seeding admin dijalankan **otomatis** saat startup.

### Tanpa Docker (Lokal)

```bash
# 1. Pastikan PostgreSQL berjalan dan set DATABASE_URL di .env
cp .env.example .env

# 2. Install dependensi
go mod download

# 3. Jalankan migrasi
make migrate-up

# 4. Build dan jalankan
make run
```

---

## Konfigurasi Lingkungan

Salin `.env.example` menjadi `.env` dan sesuaikan:

```dotenv
DATABASE_URL=postgres://tps_user:tps_pass@localhost:5432/tps_pkb?sslmode=disable
JWT_SECRET=ganti-dengan-secret-minimal-32-karakter
JWT_TTL_HOURS=24
SERVER_PORT=8080
APP_ENV=development
```

| Variabel        | Default | Keterangan                                |
|-----------------|---------|-------------------------------------------|
| `DATABASE_URL`  | —       | Connection string PostgreSQL              |
| `JWT_SECRET`    | —       | Secret HS256 (minimal 32 karakter di prod)|
| `JWT_TTL_HOURS` | `24`    | TTL token dalam jam                       |
| `SERVER_PORT`   | `8080`  | Port HTTP server                          |
| `APP_ENV`       | `development` | `development` atau `production`    |

---

## Endpoint API

Base URL: `http://localhost:8080`

| Metode | Path                          | Peran yang Diizinkan               | Keterangan                        |
|--------|-------------------------------|------------------------------------|-----------------------------------|
| GET    | `/health`                     | Publik                             | Health check                      |
| POST   | `/api/v1/auth/login`          | Publik                             | Login, dapatkan JWT               |
| GET    | `/api/v1/wajib-pajak`         | ADMIN, PETUGAS, WAJIB_PAJAK (self) | Daftar wajib pajak                |
| POST   | `/api/v1/wajib-pajak`         | ADMIN, PETUGAS                     | Buat wajib pajak baru             |
| GET    | `/api/v1/wajib-pajak/{id}`    | ADMIN, PETUGAS, WAJIB_PAJAK (self) | Detail wajib pajak                |
| PUT    | `/api/v1/wajib-pajak/{id}`    | ADMIN, PETUGAS                     | Update wajib pajak                |
| GET    | `/api/v1/kewajiban-pajak`     | ADMIN, PETUGAS, WAJIB_PAJAK (self) | Daftar kewajiban pajak            |
| POST   | `/api/v1/kewajiban-pajak`     | ADMIN                              | Buat kewajiban pajak              |
| POST   | `/api/v1/pembayaran`          | ADMIN, PETUGAS                     | Catat pembayaran + hitung denda   |
| GET    | `/api/v1/denda/{id}`          | ADMIN, PETUGAS, WAJIB_PAJAK        | Kalkulasi denda (read-only)       |
| GET    | `/api/v1/laporan`             | ADMIN, PETUGAS                     | Laporan rekapitulasi periode      |

### Contoh: Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Pretest@2025"}'
```

Response:
```json
{
  "success": true,
  "message": "login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

### Contoh: Buat Wajib Pajak

```bash
curl -X POST http://localhost:8080/api/v1/wajib-pajak \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Budi Santoso",
    "jenis": "INDIVIDU",
    "nik": "3201234567890001",
    "alamat": "Jl. Merdeka No. 1, Jakarta"
  }'
```

### Format Response Standar

**Sukses:**
```json
{ "success": true, "message": "...", "data": {...}, "meta": {...} }
```

**Gagal:**
```json
{
  "success": false,
  "message": "...",
  "error": { "code": "VALIDATION_ERROR", "details": [{"field": "nik", "message": "..."}] }
}
```

---

## Dokumentasi API (Swagger)

Setelah server berjalan, buka:

- **Swagger UI:** [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- **Raw OpenAPI YAML:** [http://localhost:8080/docs/openapi.yaml](http://localhost:8080/docs/openapi.yaml)

---

## Aturan Denda

| Jenis Denda       | Besaran                    | Kondisi                                         |
|-------------------|----------------------------|-------------------------------------------------|
| Denda Telat Bayar | 2% × pokok pajak           | `tanggal_bayar` > `periode_final`              |
| Denda Kurang Bayar| 1% × selisih kekurangan    | `total_dibayar + jumlah_bayar` < `pokok_pajak` |
| Denda Gabungan    | Keduanya dijumlahkan       | Terlambat sekaligus kurang bayar               |

Denda dihitung dan disimpan secara **atomik** bersama transaksi pembayaran dalam satu database transaction.

---

## Peran & Akses

| Peran        | Hak Akses                                                      |
|--------------|----------------------------------------------------------------|
| `ADMIN`      | Akses penuh ke semua endpoint                                  |
| `PETUGAS`    | Buat wajib pajak, catat pembayaran, lihat laporan              |
| `WAJIB_PAJAK`| Hanya baca data milik sendiri (difilter oleh `wajib_pajak_id` di JWT) |

### Kredensial Default

| Username | Password      | Peran |
|----------|---------------|-------|
| `admin`  | `Pretest@2025`| ADMIN |

---

## Menjalankan Tes

```bash
# Semua unit test
make test

# Dengan coverage report
make test-cover

# Verbose
go test ./internal/... -v
```

Strategi pengujian:

| Layer         | Target uji              | Dependensi di-mock                   |
|---------------|-------------------------|--------------------------------------|
| Handler       | HTTP request/response   | `port.XxxService`                    |
| Service       | Logika bisnis           | `port.XxxRepository`, `AuditService` |
| DendaService  | Kalkulasi denda (murni) | Tidak ada (pure function)            |

---

## Perintah Pengembang

```bash
make build        # Compile binary ke ./bin/server
make run          # Build dan jalankan server
make test         # Jalankan semua unit test
make test-cover   # Unit test + laporan coverage HTML
make migrate-up   # Jalankan migrasi database
make migrate-down # Rollback satu migrasi
make generate     # Regenerasi mock (setelah ubah interface di port/)
make lint         # Jalankan golangci-lint
```

---

## Audit Trail

Setiap mutasi data (CREATE/UPDATE) dicatat di tabel `audit_logs`:

- **Siapa** yang melakukan perubahan (`user_id`, `username`)
- **Dari mana** (`ip_address`, `user_agent`)
- **Apa yang berubah** (`old_data` → `new_data` dalam JSONB)
- **Kapan** (`created_at`)

Tabel ini **append-only** — role database aplikasi tidak memiliki izin UPDATE/DELETE pada tabel ini.
