# Build Stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git make
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
    -ldflags="-w -s" \
    -o /build/valkey-mcp-server \
    ./cmd/valkey-mcp-server

# Runtime Stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates
RUN addgroup -g 1000 valkey && \
    adduser -D -u 1000 -G valkey valkey
WORKDIR /app
COPY --from=builder /build/valkey-mcp-server /app/valkey-mcp-server
USER valkey
ENV VALKEY_URL=valkey://localhost:6379 \
    VALKEY_PASSWORD="" \
    VALKEY_DB=0
ENTRYPOINT ["/app/valkey-mcp-server"]
