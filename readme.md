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

## 🛠️ Tech Stack

- **Language**: Go 1.24+
- **RPC Framework**: ConnectRPC
- **Protocol Buffers**: buf CLI for code generation
- **TLS**: mkcert for local certificate management
- **HTTP/3**: QUIC protocol support
- **State Management**: Custom in-memory state manager
- **Observability**: Structured logging and health monitoring

## 📋 Prerequisites

Before running this service, make sure you have:

- Go 1.21 or later
- [buf CLI](https://buf.build/docs/installation)
- [mkcert](https://github.com/FiloSottile/mkcert) for local TLS certificates

## 🚀 Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/soundphilosopher/basic-grpc-service-go
cd basic-grpc-service-go
```

### 2. Generate TLS Certificates

```bash
# Install mkcert (if not already installed)
# On macOS: brew install mkcert
# On Linux: see mkcert installation docs

# Create local CA
mkcert -install

# Generate certificates for localhost
mkdir -p certs
mkcert -key-file certs/local.key -cert-file certs/local.crt localhost 127.0.0.1 ::1
```

### 3. Generate Protocol Buffer Code

```bash
buf generate
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Service

```bash
go run main.go
```

The service will start on `https://127.0.0.1:8443` with both HTTP/2 and HTTP/3 support.

## 🔌 API Endpoints

### Hello
Simple greeting service that demonstrates basic request/response pattern.

**Request:**
```json
{
  "message": "World"
}
```

**Response:** CloudEvent containing greeting message.

### Talk
Interactive bidirectional streaming for real-time conversations.

**Usage:** Send messages and receive contextual responses in real-time.

### Background
Demonstrates asynchronous processing with state management and fan-out/fan-in patterns.

**Features:**
- State tracking (PROCESS → COMPLETE)
- Concurrent external service calls
- Response aggregation
- Real-time status updates via streaming

## 🧪 Testing the Service

### Using grpcurl

```bash
# Health check
grpcurl 127.0.0.1:8443 grpc.health.v1.Health/Check

# Hello endpoint
grpcurl -d '{"message":"World"}' \
  127.0.0.1:8443 basic.v1.BasicService/Hello

# ELIZA chatbot
cat <<EOM | grpcurl -d @ 127.0.0.1:8443 basic.v1.BasicService/Talk
{
  "message": "Hello"
}
{
  "message": "How are you?"
}
{
  "message": "Bye."
}
EOM

# Background processing
grpcurl '{"processes":5}' \
  127.0.0.1:8443 basic.v1.BasicService/Background
```

### Using curl (HTTP/2)

```bash
# Hello endpoint via HTTP/2
curl -vv --http2 -X POST https://127.0.0.1:8443/basic.v1.BasicService/Hello \
  -H "Content-Type: application/json" \
  -d '{"message":"World"}'

# Hello endpoint via HTTP/3
curl -vv --http3-only -X POST https://127.0.0.1:8443/basic.v1.BasicService/Hello \
  -H "Content-Type: application/json" \
  -d '{"message":"World"}'
```

### Service Reflection

```bash
# List available services
grpcurl 127.0.0.1:8443 list

# Describe a service
grpcurl 127.0.0.1:8443 describe basic.v1.BasicService

# List methods of service
grpcurl 127.0.0.1:8443 list basic.v1.BasicService
```

## 🏗️ Project Structure

```
├── certs/              # TLS certificates
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
mkcert -key-file certs/local.key -cert-file certs/local.crt localhost 127.0.0.1 ::1
```

### Port Already in Use
If port 8443 is busy, modify the address in `main.go`:
```go
func getServerAddress() string {
    return "127.0.0.1:9000" // Change port as needed
}
```

### Protocol Buffer Generation Fails
Ensure buf CLI is properly installed and run:
```bash
buf mod update
buf generate
```

## 🔗 Useful Links

- [ConnectRPC Documentation](https://connectrpc.com/docs/)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [buf CLI Reference](https://buf.build/docs/)
- [gRPC Health Checking](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
- [HTTP/3 Specification](https://tools.ietf.org/html/rfc9114)

---

**Built with ❤️ and Go
