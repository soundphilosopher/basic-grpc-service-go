# Multi-stage Docker build for a Go gRPC service
# Stage 1: Build stage using official Go image with Alpine Linux
FROM golang:1.24-alpine3.22 AS build

# Set the working directory inside the build container
WORKDIR /app

# Copy Go module files first to leverage Docker layer caching
# This allows dependency downloads to be cached if go.mod/go.sum haven't changed
COPY ./go.* .

# Download and verify Go module dependencies
# go mod download: downloads modules to local cache
# go mod verify: verifies dependencies have expected cryptographic hashes
RUN go mod download && go mod verify

# Copy application source code
# Copying source files after dependency download ensures better layer caching
COPY ./internal ./internal
COPY ./sdk ./sdk
COPY ./main.go .

# Build the Go application
# -v: verbose output to show compilation progress
# -o: specify output binary name and location
RUN go build -v -o /usr/local/bin/grpc-server .

# Stage 2: Runtime stage using minimal Alpine Linux image
FROM alpine:3.22

# Set working directory in the runtime container
WORKDIR /app

# Copy the compiled binary from the build stage
# This creates a much smaller final image containing only the binary
COPY --from=build /usr/local/bin/grpc-server /usr/local/bin/grpc-server

# Set the default command to run when the container starts
CMD [ "grpc-server", "-server-addr", "0.0.0.0:8443" ]
