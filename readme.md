# 🚀 Basic gRPC Service

[![Go Version](https://img.shields.io/github/go-mod/go-version/soundphilosopher/basic-grpc-service-go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/soundphilosopher/basic-grpc-service-go)](https://goreportcard.com/report/github.com/soundphilosopher/basic-grpc-service-go)
[![ConnectRPC](https://img.shields.io/badge/ConnectRPC-v1.18.1-blue.svg)](https://connectrpc.com/)
[![HTTP/3](https://img.shields.io/badge/HTTP%2F3-supported-green.svg)](https://tools.ietf.org/html/rfc9114)

A modern, high-performance gRPC service built with Go, featuring ConnectRPC, dual HTTP/2 and HTTP/3 support, and comprehensive observability.

## ✨ Features

- **Modern gRPC**: Built with [ConnectRPC](https://connectrpc.com/) for better developer experience
- **Dual Protocol Support**: HTTP/2 and HTTP/3 on the same port
- **CloudEvents Integration**: All responses wrapped in CloudEvents for better event-driven architecture
- **TLS Security**: Local development certificates with mkcert
- **Health Checks**: Built-in gRPC health checking
- **Service Reflection**: Automatic service discovery and introspection
- **Streaming Support**: Bidirectional streaming capabilities
- **Background Processing**: Asynchronous task processing with state management
- **Fan-out/Fan-in Pattern**: Demonstrates concurrent service calls and response aggregation
- **Docker Support**: Multi-stage Docker build for optimized container deployment
- **Configurable Address**: Command-line flag support for server address configuration

## 🛠️ Tech Stack

- **Language**: Go 1.24+
- **RPC Framework**: ConnectRPC
- **Protocol Buffers**: buf CLI for code generation
- **TLS**: mkcert for local certificate management
- **HTTP/3**: QUIC protocol support
- **State Management**: Custom in-memory state manager
- **Observability**: Structured logging and health monitoring
- **Containerization**: Docker with Alpine Linux base

## 📋 Prerequisites

Before running this service, make sure you have:

- Go 1.21 or later
- [buf CLI](https://buf.build/docs/installation)
- [mkcert](https://github.com/FiloSottile/mkcert) for local TLS certificates
- Docker (optional, for containerized deployment)

## 🚀 Quick Start

### Option 1: Native Go Development

#### 1. Clone and Setup

```bash
git clone https://github.com/soundphilosopher/basic-grpc-service-go
cd basic-grpc-service-go
```

#### 2. Generate TLS Certificates

```bash
# Install mkcert (if not already installed)
# On macOS: brew install mkcert
# On Linux: see mkcert installation docs

# Create local CA
mkcert -install

# Generate certificates for localhost
mkdir -p certs
mkcert -key-file certs/local.key -cert-file certs/local.crt localhost 127.0.0.1 0.0.0.0 ::1
```

#### 3. Generate Protocol Buffer Code

```bash
buf generate
```

#### 4. Install Dependencies

```bash
go mod tidy
```

#### 5. Run the Service

```bash
# Run with default address (127.0.0.1:8443)
go run main.go

# Run with custom address
go run main.go -server-addr "0.0.0.0:9090"
```

### Option 2: Docker Deployment

#### 1. Generate TLS Certificates

```bash
# Create certificates directory
mkdir -p certs

# Generate certificates for localhost and Docker environments
mkcert -key-file certs/local.key -cert-file certs/local.crt localhost 127.0.0.1 0.0.0.0 ::1
```

#### 2. Build Docker Image

```bash
# Build the Docker image
docker build -f docker/Dockerfile -t basic-grpc-service:0.1.0 .
```

#### 3. Run Docker Container

```bash
# Run with default address (127.0.0.1:8443)
docker run -d \
  --name grpc-service \
  -p 8443:8443/tcp \
  -p 8443:8443/udp \
  -v $(pwd)/certs:/app/certs:ro \
  basic-grpc-service:0.1.0

# Run with custom address (bind to all interfaces)
docker run -d \
  --name grpc-service \
  -p 9090:9090/tcp \
  -p 9090:9090/udp \
  -v $(pwd)/certs:/app/certs:ro \
  basic-grpc-service:0.1.0 grpc-server -server-addr "0.0.0.0:9090"
```

#### 4. Verify Container Status

```bash
# Check container logs
docker logs grpc-service

# Test the service
curl -k --http2 -X POST https://localhost:8443/basic.v1.BasicService/Hello \
  -H "Content-Type: application/json" \
  -d '{"message":"World"}'
```

The service will start with both HTTP/2 and HTTP/3 support.

## ⚙️ Configuration

### Command Line Flags

- **`-server-addr`**: Server bind address (default: `127.0.0.1:8443`)

```bash
# Examples
./grpc-server -server-addr "0.0.0.0:8080"    # Bind to all interfaces on port 8080
./grpc-server -server-addr "localhost:9443"   # Bind to localhost on port 9443
./grpc-server -h                               # Show help with available flags
```
```

```basic-grpc-service-go/readme.md#L185-235
## 🏗️ Project Structure

```
├── certs/              # TLS certificates
├── docker/             # Docker configuration
│   └── Dockerfile     # Multi-stage Docker build
├── examples/           # Usage examples and demos
├── internal/           # Private application code
│   ├── talk/          # Conversation logic
│   └── utils/         # Utility functions
├── proto/             # Protocol buffer definitions
│   ├── basic/         # Service definitions
│   └── io/            # CloudEvents definitions
├── sdk/               # Generated gRPC code
├── buf.gen.yaml       # Buf code generation config
├── buf.yaml           # Buf project config
├── go.mod             # Go dependencies
└── main.go            # Application entry point
```

## 🔧 Development

### Regenerate Protocol Buffers

```bash
buf generate --clean
```

### Update Dependencies

```bash
go mod tidy
go mod verify
```

### Format and Lint

```bash
go fmt ./...
go vet ./...
```

### Docker Development

```bash
# Rebuild Docker image
docker build -f docker/Dockerfile -t basic-grpc-service:0.1.0 .

# Run in development mode with live logs
docker run --rm \
  -p 8443:8443/tcp \
  -p 8443:8443/udp \
  -v $(pwd)/certs:/app/certs:ro \
  basic-grpc-service:0.1.0

# Clean up containers and images
docker rm -f grpc-service
docker rmi basic-grpc-service
```

## 🚦 Health Monitoring

The service includes comprehensive health checking:

- **Health Check Endpoint**: Standard gRPC health checking
- **Service Reflection**: Automatic API documentation
- **TLS Status**: Secure connections monitoring
- **State Management**: Background task status tracking

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## 🆘 Troubleshooting

### Certificate Issues
If you encounter TLS certificate errors:
```bash
# Regenerate certificates
rm -rf certs/
mkdir certs
mkcert -key-file certs/local.key -cert-file certs/local.crt localhost 127.0.0.1 0.0.0.0 ::1
```

### Port Configuration
Change the server address using the command-line flag:
```bash
# Native Go
go run main.go -server-addr "127.0.0.1:9000"

# Docker
docker run -p 9000:9000/tcp -p 9000:9000/udp \
  -v $(pwd)/certs:/app/certs:ro \
  basic-grpc-service -server-addr "0.0.0.0:9000"
```

### Docker Issues

**Container won't start:**
```bash
# Check logs
docker logs grpc-service

# Verify certificates are mounted
docker run --rm -v $(pwd)/certs:/app/certs:ro alpine ls -la /app/certs
```

**Port binding conflicts:**
```bash
# Use different ports
docker run -p 9090:8443/tcp -p 9090:8443/udp basic-grpc-service
```

### Protocol Buffer Generation Fails
Ensure buf CLI is properly installed and run:
```bash
buf mod update
buf generate
```

## 🐳 Docker Details

The Docker image uses a multi-stage build:

1. **Build Stage**: Uses `golang:1.24-alpine3.22` to compile the Go application
2. **Runtime Stage**: Uses minimal `alpine:3.22` for a small, secure final image

**Benefits:**
- Small image size (~15MB final image)
- Security-focused with minimal attack surface
- Optimized for production deployment
- Efficient layer caching for fast rebuilds

## 🔗 Useful Links

- [ConnectRPC Documentation](https://connectrpc.com/docs/)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [buf CLI Reference](https://buf.build/docs/)
- [gRPC Health Checking](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
- [HTTP/3 Specification](https://tools.ietf.org/html/rfc9114)
- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)

---

**Built with ❤️ and Go
