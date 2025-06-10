# --- STAGE 1: Build the Go binary ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install Git (needed for some Go modules)
RUN apk add --no-cache git

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o server .

# --- STAGE 2: Create a minimal image to run ---
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Expose the port your Go app uses
EXPOSE 8080

# Run the app
CMD ["./server"]
