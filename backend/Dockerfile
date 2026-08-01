# Stage 1: Build the Go application
FROM golang:1.26-alpine AS builder

# Install SSL ca certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy dependency manifests first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build static binary with optimizations (-s -w strips debug symbol table)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/cloudstorex-server ./cmd/server

# Stage 2: Minimal Distroless non-root runtime (Refinement 5)
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

# Copy timezone data and CA certs from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /

# Copy the static binary
COPY --from=builder /app/cloudstorex-server /cloudstorex-server

# Ensure non-root execution (UID 65532 is standard nonroot in distroless)
USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/cloudstorex-server"]
