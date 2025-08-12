# ---------- Builder Stage ----------
FROM golang:1.24.4-bookworm AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y --no-install-recommends \
	librdkafka-dev \
	ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Build with CGO enabled
RUN mkdir -p bin && \
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./bin/metrics-consumer ./cmd/main.go

# ---------- Runtime Stage ----------
FROM debian:bookworm-slim

WORKDIR /app/

RUN apt-get update && apt-get install -y --no-install-recommends \
	librdkafka1 \
	ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/bin/metrics-consumer .
COPY --from=builder /build/config.yml .

ENTRYPOINT ["./metrics-consumer"]
