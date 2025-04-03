# Build stage
FROM golang:1.20-alpine AS builder

# Install required tools and dependencies
RUN apk add --no-cache git protobuf-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /gauth-extractor ./cmd/extractor

# Final stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /gauth-extractor /usr/local/bin/gauth-extractor

# Create a non-root user to run the application
RUN adduser -D -h /home/appuser appuser
USER appuser
WORKDIR /home/appuser

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/gauth-extractor", "-i"]