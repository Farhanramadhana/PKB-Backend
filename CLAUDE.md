# TPS-PKB — Tax Processing System: Pajak Kendaraan Bermotor

Samsat motor vehicle tax (PKB) REST API backend. Take-home test.

## Stack

- **Language:** Go 1.22 (stdlib `net/http` — no framework)
- **Database:** PostgreSQL 16 via `pgx/v5` + pgxpool
- **Auth:** JWT HS256 (`golang-jwt/jwt/v5`), bcrypt password hashing (cost 12)
- **Money:** `shopspring/decimal` — never `float64` for currency
- **Validation:** `go-playground/validator/v10` with custom NIK/NPWP/NIB validators
- **SQL:** Raw `pgx` queries + `squirrel` for dynamic filter queries
- **Migrations:** `golang-migrate/migrate/v4` — runs automatically on startup
- **Mocks:** `go.uber.org/mock/mockgen` — run `make generate` to regenerate
- **Config:** `viper` — reads `.env` file or environment variables

## Architecture: Hexagonal (Ports & Adapters)

```
cmd/server/main.go          ← wire-up only; all deps injected here
internal/domain/
  entity/                   ← pure Go structs, no imports from outer layers
  dto/                      ← shared request/filter types used in port interfaces
  port/                     ← interfaces: service.go, repository.go, audit.go, token_provider.go
  service/                  ← business logic; depends only on port interfaces
internal/adapter/
  http/                     ← driving adapter (handlers, middleware, router, response)
  postgres/                 ← driven adapter (repository implementations)
internal/infrastructure/
  config/                   ← Viper config binding
  database/                 ← pgxpool setup + migration runner
  token/                    ← JWT TokenProvider implementation
internal/mock/              ← generated mocks (commit these; regenerate with make generate)
migrations/                 ← numbered .up.sql / .down.sql files
```

**The one rule:** `internal/domain/` never imports from `adapter/` or `infrastructure/`.

## API Versioning

All routes: `/api/v1/...`  
Path params accessed via `r.PathValue("id")` (Go 1.22).

## Standard API Response

Every handler calls exactly one of:
- `response.WriteSuccess(w, statusCode, message, data)`
- `response.WriteSuccessPaginated(w, statusCode, message, data, meta)`
- `response.WriteError(w, statusCode, code, message, details...)`

Success envelope:
```json
{ "success": true, "message": "...", "data": {...}, "meta": {...} }
```
Error envelope:
```json
{ "success": false, "message": "...", "error": { "code": "VALIDATION_ERROR", "details": [...] } }
```

## Domain Language (Glossary)

| Term | Meaning |
|---|---|
| Wajib Pajak | Taxpayer — individual (NIK) or corporate (NPWP/NIB) |
| Kendaraan | Motor vehicle registered under a wajib pajak |
| Kewajiban Pajak | Annual tax obligation per vehicle |
| Pembayaran | A payment transaction against a kewajiban |
| Denda Telat Bayar | Late payment fine: 2% × pokok pajak |
| Denda Kurang Bayar | Underpayment fine: 1% × (pokok − dibayar) |
| Denda Gabungan | Both fines applied when late AND underpaid |
| Laporan | Period payment report: total kewajiban, dibayar, denda, sisa |
| Pokok Pajak | Base tax principal amount |

## Roles & Access

| Role | Access |
|---|---|
| ADMIN | Full access to all endpoints |
| PETUGAS | Create wajib pajak, record payments, view reports |
| WAJIB_PAJAK | Read own data only (filtered by wajib_pajak_id in JWT) |

## Test Credentials

```
username: admin
password: Pretest@2025
```

## Validation Rules

| Field | Rule |
|---|---|
| NIK | exactly 16 numeric digits |
| NPWP | exactly 15 numeric digits |
| NIB | non-empty alphanumeric |
| jumlah_bayar | > 0 |
| tanggal_bayar | ≤ today (not future) |
| wajib_pajak | must be active (is_active = true) when recording payment |

## Audit Trail

Every mutation (CREATE/UPDATE) is written to `audit_logs` by the service layer.  
The service reads `user_id`, `username`, `ip_address`, `user_agent` from `context.Context` (injected by `AuthMiddleware`).  
`audit_logs` is append-only — the DB role has INSERT but no UPDATE/DELETE on it.

## Developer Commands

```bash
# Start all services
docker-compose up --build

# Run migrations only
make migrate-up

# Generate mocks (after changing port interfaces)
make generate

# Run tests
make test

# Run with coverage
make test-cover

# Build binary
make build
```

## Adding a New Endpoint

1. Add method to the appropriate interface in `internal/domain/port/service.go`
2. Implement the method in `internal/domain/service/`
3. Add method to repository interface in `internal/domain/port/repository.go` (if needed)
4. Implement in `internal/adapter/postgres/`
5. Add handler method in `internal/adapter/http/handler/`
6. Register route in `internal/adapter/http/router.go`
7. Run `make generate` to update mocks
8. Write unit tests for handler + service

## Environment Variables

See `.env.example`. Key vars:

| Var | Default | Description |
|---|---|---|
| `DATABASE_URL` | — | PostgreSQL connection string |
| `JWT_SECRET` | — | HS256 signing secret (min 32 chars in prod) |
| `JWT_TTL_HOURS` | `24` | Token TTL in hours |
| `SERVER_PORT` | `8080` | HTTP listen port |
| `APP_ENV` | `development` | `development` or `production` |
