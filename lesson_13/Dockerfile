# Build stage
FROM golang:1.23.2-alpine AS build

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY internal/ internal/
COPY cmd/ cmd/

# Build the server
# Disable CGO for a fully static binary
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /tmp/documentstore-server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS if needed
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from build stage
COPY --from=build /tmp/documentstore-server /app/server

# Create data directory for potential file storage
RUN mkdir -p /data

# Expose the server port
EXPOSE 8080

# Run the server
CMD ["/app/server"]
