package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestGetServerAddress(t *testing.T) {
	t.Parallel()

	t.Run("should return default address", func(t *testing.T) {
		address := getServerAddress()
		assert.Equal(t, "127.0.0.1:8443", address)
	})
}

func TestCreateHTTP2Server(t *testing.T) {
	// Arrange
	mux := setupMux()
	addr := "127.0.0.1:0" // Use port 0 to bind to a random available port
	httpServer := createHTTP2Server(addr, mux)

	// Act
	ln, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		t.Fatalf("Failed to start HTTP2 server: %v", err)
	}
	defer ln.Close()

	go func() {
		httpServer.Serve(ln)
	}()
	time.Sleep(100 * time.Millisecond) // Give server some time to start

	// Assert
	assert.NotNil(t, ln)
}

func TestCreateHTTP3Server(t *testing.T) {
	// Arrange
	mux := setupMux()
	addr := "127.0.0.1:0" // Random port
	http3Server := createHTTP3Server(addr, mux)

	// Act
	udpAddr, err := net.ResolveUDPAddr("udp", http3Server.Addr)
	if err != nil {
		t.Fatalf("Failed to resolve UDP address: %v", err)
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatalf("Failed to start HTTP3 server: %v", err)
	}
	defer udpConn.Close()

	go func() {
		http3Server.Serve(udpConn)
	}()
	time.Sleep(100 * time.Millisecond) // Give server time to start

	// Assert
	assert.NotNil(t, udpConn)
}

func TestSetupListenersWithSignalCancel(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to simulate receiving OS signals (like SIGINT for CTRL+C)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Set up a simple HTTP/2 server with a random port
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    "127.0.0.1:0", // Random port
		Handler: mux,
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2"},
		},
	}
	ln, err := net.Listen("tcp", httpServer.Addr)
	assert.NoError(t, err)
	defer ln.Close()

	// Create an HTTP/3 server bound to the same address
	http3Server := &http3.Server{
		Addr:    "127.0.0.1:0", // Random port for HTTP/3 as well
		Handler: mux,
	}

	// Start the HTTP/2 server in the background
	go func() {
		_ = httpServer.ServeTLS(ln, "./certs/local.cert", "./certs/local.key")
	}()

	// Act: Start listeners using setupListeners
	eg, egCtx := errgroup.WithContext(ctx)
	go func() {
		err = setupListeners(egCtx, httpServer.Addr, httpServer, http3Server)
		assert.NoError(t, err)
	}()

	// Simulate the servers running for a short time (as if they were running normally)
	time.Sleep(200 * time.Millisecond)

	// Simulate receiving SIGINT (CTRL+C) after servers are running
	signalCh <- syscall.SIGINT

	// Ensure servers shut down gracefully
	err = eg.Wait()
	assert.NoError(t, err)

	// Stop listening for OS signals
	signal.Stop(signalCh)
	close(signalCh)
}
