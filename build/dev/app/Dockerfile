FROM golang:1.24-alpine AS builder

WORKDIR /app

# Кэшируем зависимости Go
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go env -w GOMODCACHE=/go/pkg/mod

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app ./cmd/app

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/bin/app .
COPY --from=builder /app/build/dev/.env .

CMD ["./app", "-broker_addr=kafka0:29092", "-addr=0.0.0.0:8081"]
