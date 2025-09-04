// Package main implements a gRPC service supporting both HTTP/2 and HTTP/3 protocols
// using Connect. The service provides basic operations with health checks and reflection.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"github.com/quic-go/quic-go/http3"
	"github.com/soundphilosopher/basic-grpc-service-go/internal"
	"github.com/soundphilosopher/basic-grpc-service-go/sdk/basic/v1/basicV1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
)

// main starts both HTTP/2 and HTTP/3 servers concurrently on the same address.
// Both servers serve the same gRPC services with TLS enabled.
func main() {
	mux := setupMux()
	addr := getServerAddress()

	httpServer := createHTTP2Server(addr, mux)
	defer httpServer.Close()

	http3Server := createHTTP3Server(addr, mux)
	defer http3Server.Close()

	if err := setupListeners(context.Background(), addr, &httpServer, &http3Server); err != nil {
		log.Fatalf("failed to setup listeners: %v", err)
	}
}

// setupMux configures the HTTP multiplexer with gRPC services, health checks,
// and reflection handlers. All handlers use 1KB minimum compression.
func setupMux() *http.ServeMux {
	compress1KB := connect.WithCompressMinBytes(1024)
	mux := http.NewServeMux()

	// Register core business service
	mux.Handle(basicV1connect.NewBasicServiceHandler(internal.NewBasicServiceV1(), compress1KB))

	// Register health and reflection services
	checkServices := []string{
		basicV1connect.BasicServiceName,
	}
	mux.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker(checkServices...), compress1KB))
	mux.Handle(grpcreflect.NewHandlerV1(grpcreflect.NewStaticReflector(checkServices...), compress1KB))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(checkServices...), compress1KB))

	return mux
}

// getServerAddress parses command line flags and returns the server bind address.
// If -server-addr flag is not provided, defaults to "127.0.0.1:8443"
func getServerAddress() string {
	serverAddr := flag.String("server-addr", "127.0.0.1:8443", "server address to bind to")
	flag.Parse()

	return *serverAddr
}

// createHTTP2Server creates an HTTP/2 server with h2c support and reasonable timeouts.
func createHTTP2Server(addr string, handler http.Handler) http.Server {
	return http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(handler, &http2.Server{}),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024,
	}
}

// createHTTP3Server creates an HTTP/3 server using QUIC protocol.
func createHTTP3Server(addr string, handler http.Handler) http3.Server {
	return http3.Server{
		Addr:    addr,
		Handler: handler,
	}
}

// setupListeners starts both HTTP/2 and HTTP/3 servers concurrently.
// Returns an error if either server fails to start.
func setupListeners(ctx context.Context, addr string, httpServer *http.Server, http3Server *http3.Server) error {
	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		log.Printf("Start HTTP over TCP server on %s ...", addr)
		return httpServer.ListenAndServeTLS("./certs/local.crt", "./certs/local.key")
	})

	eg.Go(func() error {
		log.Printf("Start HTTP over UDP server on %s ...", addr)
		return http3Server.ListenAndServeTLS("./certs/local.crt", "./certs/local.key")
	})

	return eg.Wait()
}
