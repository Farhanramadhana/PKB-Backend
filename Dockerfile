FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /tps-pkb ./cmd/server

# ---

FROM alpine:3.19

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /tps-pkb .
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 8080

CMD ["./tps-pkb"]
