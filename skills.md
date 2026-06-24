# Engineer Skills Profile — TPS-PKB Project

## Role
Senior Backend Engineer — sole implementor of this take-home system.

## Core Skills Applied

### Go (Golang)
- Go 1.22 idiomatic patterns: stdlib `net/http` routing with method+path patterns, `context` propagation, interface-driven design
- Unexported concrete types with exported constructor returning interface — forces callers to code against the abstraction
- `shopspring/decimal` for monetary arithmetic — never `float64`
- `pgx/v5` native PostgreSQL driver (no ORM); raw SQL with `squirrel` for dynamic queries
- `golang-migrate` for schema versioning; migrations run automatically at startup
- `log/slog` structured JSON logging (stdlib, Go 1.21+)
- `go-playground/validator/v10` with custom validation functions for NIK (16 digits), NPWP (15 digits), NIB

### Architecture
- **Hexagonal (Ports & Adapters):** domain layer has zero infrastructure imports; adapters implement domain interfaces
- Service interfaces defined in `internal/domain/port/service.go` — handlers depend on these, never on concrete service types
- Repository interfaces in `internal/domain/port/repository.go` — services depend on these, enabling full unit-test isolation
- Pure domain service for fine calculation (`DendaService.Hitung`) — no I/O, inputs → `[]entity.Denda`; trivially testable
- Clean separation: `entity` (structs), `dto` (in/out data shapes), `port` (interfaces), `service` (logic), `adapter` (I/O)

### PostgreSQL
- Schema design with ENUMs, CHECK constraints (enforce NIK/NPWP/NIB business rules at DB level), partial indexes
- Atomic multi-table writes via `pgx.Tx` (`SaveWithDenda`: pembayaran + denda + kewajiban status in one transaction)
- Connection pooling with `pgxpool`
- Append-only audit log table — DB role has INSERT-only permission

### API Design
- URL versioning: `/api/v1/...` — explicit, proxy-friendly, forward-compatible
- Standard response envelope: `{ success, message, data, meta?, error? }` — all handlers use `response.WriteSuccess` / `response.WriteError`
- HTTP error codes mapped to domain error codes: `VALIDATION_ERROR`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`, `CONFLICT`, `UNPROCESSABLE`, `INTERNAL_ERROR`
- JWT HS256 authentication; RBAC middleware with `RequireRole` + `SelfOnly` for WAJIB_PAJAK data isolation

### Security
- bcrypt cost 12 for password hashing
- JWT claims carry `user_id`, `role`, `wajib_pajak_id` (for WAJIB_PAJAK role ownership checks)
- Input validation at HTTP boundary; domain validates business rules separately
- Audit trail: every mutation logged with user identity, IP, user-agent, before/after state (JSONB)

### Testing
- **Handler tests:** `httptest.NewRecorder` + mock service — test HTTP contract, status codes, response shape
- **Service tests:** mock repos + mock audit — test business logic in isolation
- **Pure function tests:** `DendaService.Hitung` — no mocks, deterministic, fast
- **Repository tests:** `testcontainers-go` real PostgreSQL — test SQL correctness
- Mock generation: `go.uber.org/mock/mockgen` via `//go:generate` directives; `make generate`
- Test helper: `withClaims(ctx, claims)` injects auth context for handler tests

### Docker
- Multi-stage Dockerfile: `golang:1.22-alpine` builder → `alpine:3.19` runtime; non-root user
- `docker-compose.yml`: app depends on DB with `condition: service_healthy`; Postgres healthcheck with `pg_isready`

## Domain Knowledge: Pajak Kendaraan Bermotor (PKB)

- **Wajib Pajak types:** Individu (NIK-based) and Badan Usaha (NPWP + NIB)
- **Annual tax cycle:** each vehicle generates one `kewajiban_pajak` per `tahun_pajak` with a `periode_final` (due date)
- **Payment statuses:** LUNAS (paid in full), KURANG_BAYAR (underpaid), LEBIH_BAYAR (overpaid)
- **Fine rules:**
  - Denda Telat Bayar: 2% × pokok pajak — triggered when `tanggal_bayar > periode_final`
  - Denda Kurang Bayar: 1% × shortfall — triggered when total paid < pokok pajak
  - Denda Gabungan: both fines applied simultaneously when late AND underpaid
- **Laporan (report):** aggregates over date range — total kewajiban, total dibayar, total denda, sisa kewajiban
- **SAMSAT roles:** ADMIN (full), PETUGAS (input + report), WAJIB_PAJAK (read own data)

## Microservice Migration Path

Current monolith ports map directly to future service boundaries:
- `port.AuthService` → auth-service (stateless JWT, extract first)
- `port.WajibPajakRepository` → wajib-pajak-service
- `port.PembayaranRepository` (with denda) → pembayaran-service (transactional unit, extract together)
- `port.LaporanRepository` → laporan-service (read-only, point at read replica)
