# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/vigilante ./cmd/vigilante/*.go

# Final stage
FROM scratch

# Copy compiled binary from builder
COPY --from=builder /app/bin/vigilante /vigilante

# Copy necessary files like the dashboard
COPY dashboard/index.html /dashboard/index.html
COPY internal/storage/schema.sql /internal/storage/schema.sql

# Expose HTTP and gRPC ports
EXPOSE 3000 50051

ENTRYPOINT ["/vigilante"]
CMD ["serve"]
