# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install swag for Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation and build the application
RUN make build

# Stage 2: Create the final image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/energywar .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./energywar"]
