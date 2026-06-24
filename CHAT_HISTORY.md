# Chat History — TPS-PKB Take-Home Test

Rekap seluruh percakapan sesi pengembangan sistem TPS-PKB bersama Claude (Senior Go Engineer).

---

## Sesi 1 — Perencanaan & Desain Sistem

### Permintaan Awal
User meminta bantuan membangun backend REST API untuk sistem **Pajak Kendaraan Bermotor (PKB)** yang digunakan oleh Samsat Indonesia.

### Keputusan Arsitektur
- **Hexagonal (Ports & Adapters) Monolith** — bukan microservices untuk MVP, tapi didesain mudah dipecah ke microservices nanti
- **Native Go `net/http`** — tanpa framework (chi, gin, echo)
- **Interface-driven design** — setiap boundary layer diekspresikan sebagai interface Go
- **PostgreSQL 16** via `pgx/v5` + pgxpool (tanpa ORM)
- **`shopspring/decimal`** untuk semua nilai uang — tidak pernah `float64`

### ERD yang Didefinisikan
6 tabel utama:
- `wajib_pajak` — individu (NIK) atau badan usaha (NPWP + NIB)
- `users` — autentikasi dengan role ADMIN/PETUGAS/WAJIB_PAJAK
- `kendaraan` — kendaraan terdaftar per wajib pajak
- `kewajiban_pajak` — kewajiban tahunan per kendaraan
- `pembayaran` — transaksi pembayaran
- `denda` — denda telat/kurang bayar

---

## Sesi 2 — Implementasi Penuh

### File yang Dibuat

#### Domain Layer (`internal/domain/`)
| File | Isi |
|---|---|
| `entity/user.go` | Struct User, enum Role |
| `entity/wajib_pajak.go` | Struct WajibPajak, enum JenisWajibPajak |
| `entity/kendaraan.go` | Struct Kendaraan |
| `entity/kewajiban.go` | Struct KewajibanPajak, enum StatusKewajiban |
| `entity/pembayaran.go` | Struct Pembayaran, enum StatusPembayaran |
| `entity/denda.go` | Struct Denda, enum JenisDenda (TELAT_BAYAR, KURANG_BAYAR) |
| `entity/laporan.go` | Struct Laporan (agregasi laporan periode) |
| `dto/errors.go` | Sentinel errors domain (ErrNotFound, ErrWajibPajakInaktif, dll) |
| `port/service.go` | Interface semua service (AuthService, WajibPajakService, dll) |
| `port/repository.go` | Interface semua repository |
| `port/audit.go` | Interface AuditService + struct AuditEntry |
| `port/token_provider.go` | Interface TokenProvider + struct Claims |

#### Service Layer (`internal/domain/service/`)
| File | Kunci Implementasi |
|---|---|
| `auth_service.go` | Login: cari user → bcrypt verify → generate JWT |
| `wajib_pajak_service.go` | CRUD + cek duplikat NIK/NPWP |
| `kewajiban_service.go` | Create kewajiban, validasi kendaraan & wajib pajak |
| `pembayaran_service.go` | Validasi bisnis → hitung denda → simpan atomik |
| `denda_service.go` | **Pure function** — kalkulasi denda tanpa I/O |
| `laporan_service.go` | Agregasi laporan dari repository |

**Aturan denda (`denda_service.go`):**
```
TELAT_BAYAR  : 2% × pokok_pajak        (jika tanggal_bayar > periode_final)
KURANG_BAYAR : 1% × selisih kekurangan (jika total_dibayar + jumlah < pokok)
GABUNGAN     : keduanya jika dua kondisi terpenuhi
```

#### HTTP Adapter (`internal/adapter/http/`)
| File | Isi |
|---|---|
| `response/response.go` | Envelope standar WriteSuccess / WriteError |
| `middleware/auth_middleware.go` | JWT validate → inject claims + IP + User-Agent ke context |
| `middleware/rbac_middleware.go` | RequireRole + SelfOnly + Chain helper |
| `handler/auth_handler.go` | POST /api/v1/auth/login |
| `handler/wajib_pajak_handler.go` | CRUD wajib pajak |
| `handler/kewajiban_handler.go` | List + Create kewajiban |
| `handler/pembayaran_handler.go` | POST /api/v1/pembayaran |
| `handler/denda_handler.go` | GET /api/v1/denda/{id} |
| `handler/laporan_handler.go` | GET /api/v1/laporan |
| `handler/health_handler.go` | GET /health (publik) |
| `handler/docs_handler.go` | GET /swagger/ + GET /docs/openapi.yaml |
| `router.go` | ServeMux + middleware chain untuk semua route |

#### PostgreSQL Adapter (`internal/adapter/postgres/`)
| File | Isi |
|---|---|
| `user_repo.go` | FindByUsername, FindByID |
| `wajib_pajak_repo.go` | CRUD + ExistsByNIK/NPWP + squirrel filter |
| `kendaraan_repo.go` | FindByID |
| `kewajiban_repo.go` | Save, Update, FindByID, FindAll (squirrel filter) |
| `pembayaran_repo.go` | **SaveWithDenda** — satu pgx.Tx untuk pembayaran+denda+update kewajiban |
| `denda_repo.go` | FindByKewajibanID, FindByPembayaranID |
| `laporan_repo.go` | GetLaporan (agregasi SQL) |
| `audit_repo.go` | INSERT audit_logs (append-only) |

#### Infrastructure (`internal/infrastructure/`)
| File | Isi |
|---|---|
| `config/config.go` | Viper binding `.env` |
| `database/postgres.go` | pgxpool setup + golang-migrate runner |
| `token/jwt_provider.go` | JWT HS256 Generate + Validate |

#### Migrations
| File | Isi |
|---|---|
| `000001_init_schema.up.sql` | 5 ENUM + 6 tabel + semua index |
| `000002_audit_logs.up.sql` | Tabel audit_logs (append-only) |
| `000003_seed_admin.up.sql` | INSERT admin (idempotent, bcrypt cost 12) |
| `000004_seed_dummy_data.up.sql` | 10 data dummy per tabel (lihat Sesi 5) |

#### Lainnya
| File | Isi |
|---|---|
| `cmd/server/main.go` | Wire-up semua dependency + graceful shutdown |
| `Dockerfile` | Multi-stage: golang:1.22-alpine → alpine:3.19 |
| `docker-compose.yml` | App + PostgreSQL dengan healthcheck |
| `Makefile` | build, run, test, migrate, generate |
| `.env.example` | Template konfigurasi |
| `go.mod` | Module: github.com/bpka/tps-pkb |
| `CLAUDE.md` | Dokumentasi proyek untuk Claude |

---

## Sesi 3 — Bug Fix: Mock Nil Panic

### Masalah
Semua method mock yang mengembalikan `error` atau `bool` menggunakan bare type assertion:
```go
// SALAH — panic ketika Return(nil) dipanggil di test
return ret[0].(error)
return ret[0].(bool), ret[1].(error)
```

### Root Cause
Ketika mock dikonfigurasi `.Return(nil)`, `ret[0]` adalah untyped nil. Bare type assertion `.(error)` pada untyped nil menyebabkan runtime panic: `interface conversion: interface is nil, not error`.

### Fix
Semua method diubah ke comma-ok pattern:
```go
// BENAR — aman untuk nil
ret0, _ := ret[0].(error)
return ret0

ret0, _ := ret[0].(bool)
ret1, _ := ret[1].(error)
return ret0, ret1
```

### File yang Diperbaiki
- `internal/mock/mock_repository.go`:
  - `MockPembayaranRepository.SaveWithDenda`
  - `MockWajibPajakRepository.Save`
  - `MockWajibPajakRepository.Update`
  - `MockWajibPajakRepository.ExistsByNIK`
  - `MockWajibPajakRepository.ExistsByNPWP`
  - `MockKewajibanRepository.Save`
  - `MockKewajibanRepository.Update`

### Hasil Test
```
ok  github.com/bpka/tps-pkb/internal/adapter/http/handler   (4 tests)
ok  github.com/bpka/tps-pkb/internal/domain/service         (11 tests)
```

**15 tests, semua PASS.**

---

## Sesi 4 — Health Check, Swagger, README

### 1. Health Check Endpoint
```
GET /health  →  200 {"status":"up","service":"tps-pkb"}
```
- File baru: `internal/adapter/http/handler/health_handler.go`
- Tidak memerlukan autentikasi

### 2. Swagger / API Documentation
- `internal/adapter/http/handler/openapi.yaml` — spec OpenAPI 3.0 lengkap (11 endpoint, semua schema, contoh request/response, error codes), di-embed ke binary via `//go:embed`
- `internal/adapter/http/handler/docs_handler.go` — dua endpoint:
  - `GET /swagger/` — Swagger UI (load dari CDN)
  - `GET /docs/openapi.yaml` — raw YAML spec

### 3. README.md
Dibuat dari awal, mencakup:
- Fitur utama, tech stack, arsitektur
- Quick start Docker + lokal
- Tabel endpoint API + contoh curl
- Aturan denda, RBAC, kredensial test
- Dev commands, audit trail design

---

## Sesi 5 — Seed Data Dummy

### Migration 000004
File: `migrations/000004_seed_dummy_data.up.sql`

**10 Wajib Pajak:**
| # | Nama | Jenis | Identifier |
|---|---|---|---|
| 1 | Budi Santoso | INDIVIDU | NIK 3171050415850001 |
| 2 | Siti Rahayu | INDIVIDU | NIK 3273106203900002 |
| 3 | Ahmad Fauzi | INDIVIDU | NIK 3578011501880003 |
| 4 | Dewi Permatasari | INDIVIDU | NIK 3404126208920004 |
| 5 | Rizky Pratama | INDIVIDU | NIK 3603010295950005 |
| 6 | Fitria Handayani | INDIVIDU | NIK 1271085708880006 |
| 7 | Hendra Kusuma | INDIVIDU | NIK 7371230185870007 |
| 8 | PT Maju Bersama Tbk | BADAN_USAHA | NPWP 012345678901230 |
| 9 | CV Sumber Rezeki | BADAN_USAHA | NPWP 023456789012340 |
| 10 | PT Trans Nusantara Logistik | BADAN_USAHA | NPWP 034567890123450 |

**10 Kendaraan (nomor polisi Indonesia):**
| # | Kendaraan | Plat | Pemilik |
|---|---|---|---|
| 1 | Toyota Avanza 1.3 G 2020 | B 1234 ABC | Budi Santoso |
| 2 | Honda Beat 110 CBS 2021 | D 5678 DEF | Siti Rahayu |
| 3 | Mitsubishi Xpander Ultimate 2019 | L 2345 GHI | Ahmad Fauzi |
| 4 | Yamaha NMAX 155 Connected 2022 | AB 6789 JKL | Dewi Permatasari |
| 5 | Suzuki Ertiga GX MT 2018 | A 3456 MNO | Rizky Pratama |
| 6 | Honda Vario 160 CBS 2023 | BK 7890 PQR | Fitria Handayani |
| 7 | Toyota Rush TRD Sportivo 2020 | DD 4567 STU | Hendra Kusuma |
| 8 | Daihatsu Xenia 1.3 R 2021 | B 8901 VWX | PT Maju Bersama |
| 9 | Isuzu Elf NLR 55 BL 2019 | L 5678 YZA | CV Sumber Rezeki |
| 10 | Hino Ranger FM 260 JD 2018 | D 9012 BCD | PT Trans Nusantara |

**10 Kewajiban Pajak (tahun 2025) + 7 Pembayaran:**
| Kewajiban | Pokok | Status | Skenario Denda |
|---|---|---|---|
| kw01 (Avanza Budi) | Rp 2.500.000 | LUNAS | Tepat waktu — tidak ada denda |
| kw02 (Beat Siti) | Rp 500.000 | LUNAS | Terlambat → TELAT_BAYAR: **Rp 10.000** |
| kw03 (Xpander Ahmad) | Rp 3.000.000 | KURANG_BAYAR | Terlambat + kurang → GABUNGAN: **Rp 70.000** |
| kw04 (NMAX Dewi) | Rp 750.000 | LUNAS | Tepat waktu — tidak ada denda |
| kw05 (Ertiga Rizky) | Rp 2.000.000 | BELUM_BAYAR | Belum bayar |
| kw06 (Vario Fitria) | Rp 600.000 | LUNAS | Tepat waktu — tidak ada denda |
| kw07 (Rush Hendra) | Rp 2.750.000 | BELUM_BAYAR | Belum bayar |
| kw08 (Xenia PT Maju) | Rp 2.200.000 | LUNAS | Terlambat → TELAT_BAYAR: **Rp 44.000** |
| kw09 (Isuzu CV Sumber) | Rp 4.500.000 | BELUM_BAYAR | Belum bayar |
| kw10 (Hino PT Trans) | Rp 5.000.000 | KURANG_BAYAR | Tepat waktu + kurang → KURANG_BAYAR: **Rp 15.000** |

**5 Denda (semua skenario tercakup):**
| Denda | Jenis | Dasar | Tarif | Jumlah |
|---|---|---|---|---|
| d01 | TELAT_BAYAR | Rp 500.000 | 2% | Rp 10.000 |
| d02 | TELAT_BAYAR | Rp 3.000.000 | 2% | Rp 60.000 |
| d03 | KURANG_BAYAR | Rp 1.000.000 | 1% | Rp 10.000 |
| d04 | TELAT_BAYAR | Rp 2.200.000 | 2% | Rp 44.000 |
| d05 | KURANG_BAYAR | Rp 1.500.000 | 1% | Rp 15.000 |

**10 Users seed (password semua: `Pretest@2025`):**
`petugas1`, `petugas2`, `wp_budi`, `wp_siti`, `wp_ahmad`, `wp_dewi`, `wp_rizky`, `wp_fitria`, `wp_hendra`, `wp_majubersama`

---

## Sesi 6 — Diagram Alur Request di README

### Diagram yang Ditambahkan

**1. Sequence Diagram — Alur Umum**
Menunjukkan alur end-to-end semua protected endpoint:
`Client → Router → Auth MW → RequireRole → SelfOnly → Handler → Service → Repository → PostgreSQL`
Termasuk semua titik kegagalan (401, 403, 400, 422, 500).

**2. Flowchart — POST /api/v1/pembayaran**
Alur paling kompleks: validasi bisnis → kalkulasi denda murni → penyimpanan atomik satu transaksi DB.

**3. ASCII Layer Diagram**
Menunjukkan arah dependency antar layer hexagonal.

---

## Ringkasan Middleware & Guard Endpoint

| Endpoint | Auth JWT | Role Guard | SelfOnly |
|---|:---:|:---:|:---:|
| `GET /health` | ✗ | ✗ | ✗ |
| `GET /swagger/`, `/docs/openapi.yaml` | ✗ | ✗ | ✗ |
| `POST /api/v1/auth/login` | ✗ | ✗ | ✗ |
| `GET /api/v1/wajib-pajak` | ✓ | ALL | ✓ |
| `POST /api/v1/wajib-pajak` | ✓ | ADMIN, PETUGAS | ✗ |
| `GET /api/v1/wajib-pajak/{id}` | ✓ | ALL | ✗ |
| `PUT /api/v1/wajib-pajak/{id}` | ✓ | ADMIN, PETUGAS | ✗ |
| `GET /api/v1/kewajiban-pajak` | ✓ | ALL | ✓ |
| `POST /api/v1/kewajiban-pajak` | ✓ | **ADMIN saja** | ✗ |
| `POST /api/v1/pembayaran` | ✓ | ADMIN, PETUGAS | ✗ |
| `GET /api/v1/laporan` | ✓ | ADMIN, PETUGAS | ✗ |
| `GET /api/v1/denda/{id}` | ✓ | ALL | ✗ |

> **Catatan gap:** `GET /api/v1/wajib-pajak/{id}` tidak memiliki `SelfOnly`, sehingga WAJIB_PAJAK yang mengetahui UUID dapat mengakses data orang lain. Perlu ditambahkan ownership check di handler level.

---

## Kredensial & Akses

| Akun | Username | Password | Role |
|---|---|---|---|
| Admin default | `admin` | `Pretest@2025` | ADMIN |
| Petugas 1 | `petugas1` | `Pretest@2025` | PETUGAS |
| Petugas 2 | `petugas2` | `Pretest@2025` | PETUGAS |
| Budi Santoso | `wp_budi` | `Pretest@2025` | WAJIB_PAJAK |
| Siti Rahayu | `wp_siti` | `Pretest@2025` | WAJIB_PAJAK |

---

## Issues & Fixes Selama Sesi

| # | Masalah | Solusi |
|---|---|---|
| 1 | Missing `pgxpool` dep (`go.sum` entry hilang) | `go get github.com/jackc/pgx/v5/pgxpool && go mod tidy` |
| 2 | `decimal.Equal` gagal compare (exponent berbeda meski nilai sama) | Gunakan `decimal.NewFromInt(X).Equal(result)` bukan `assert.Equal` |
| 3 | Mock panic: `interface is nil, not error` | Ganti bare assertion `.(error)` ke comma-ok `ret0, _ := ret[0].(error)` di semua mock method |
| 4 | IDE gopls warning: module not in workspace | `go work init .` — warning kosmetik, build & test tetap jalan normal |
| 5 | IDE error: `use of internal package not allowed` | False positive gopls — `go build ./...` berjalan bersih |

---

## Stack Lengkap

```
Language    : Go 1.22
Framework   : stdlib net/http (no framework)
Database    : PostgreSQL 16
Driver      : pgx/v5 + pgxpool
Auth        : JWT HS256 (golang-jwt/jwt/v5)
Password    : bcrypt cost 12
Money       : shopspring/decimal
Validation  : go-playground/validator/v10
SQL Builder : squirrel (dynamic filters)
Migrations  : golang-migrate/migrate/v4
Config      : viper (.env + env vars)
Mock        : go.uber.org/mock/mockgen
Test        : testify/assert + testify/require
Container   : Docker multi-stage + docker-compose
```
