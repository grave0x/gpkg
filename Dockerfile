# Multi-stage build for gpkg
# Stage 1: Build the binary
FROM golang:1.22-alpine AS builder

# Install build dependencies (SQLite requires CGO)
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    git \
    ca-certificates

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version info
ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.version=${VERSION} -s -w" \
    -trimpath \
    -o /build/gpkg \
    ./cmd/gpkg

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs \
    git \
    && addgroup -g 1000 gpkg \
    && adduser -D -u 1000 -G gpkg gpkg

# Copy binary from builder
COPY --from=builder /build/gpkg /usr/local/bin/gpkg

# Set up user and directories
USER gpkg
WORKDIR /home/gpkg

# Create default gpkg directories
RUN mkdir -p /home/gpkg/.gpkg/{bin,packages,cache,db}

# Add metadata labels
LABEL org.opencontainers.image.title="gpkg" \
      org.opencontainers.image.description="Simple package manager for GitHub releases and source builds" \
      org.opencontainers.image.vendor="grave0x" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.url="https://github.com/grave0x/gpkg" \
      org.opencontainers.image.source="https://github.com/grave0x/gpkg" \
      org.opencontainers.image.documentation="https://github.com/grave0x/gpkg/blob/main/README.md"

# Set PATH to include gpkg bin directory
ENV PATH="/home/gpkg/.gpkg/bin:${PATH}"

ENTRYPOINT ["gpkg"]
CMD ["--help"]
