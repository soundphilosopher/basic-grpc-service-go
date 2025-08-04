package main

import (
	"context"
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

func setupMux() *http.ServeMux {
	compress1KB := connect.WithCompressMinBytes(1024)
	mux := http.NewServeMux()

	mux.Handle(basicV1connect.NewBasicServiceHandler(internal.NewBasicServiceV1(), compress1KB))
	mux.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker(basicV1connect.BasicServiceName), compress1KB))
	mux.Handle(grpcreflect.NewHandlerV1(grpcreflect.NewStaticReflector(basicV1connect.BasicServiceName), compress1KB))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(grpcreflect.NewStaticReflector(basicV1connect.BasicServiceName), compress1KB))

	return mux
}

func getServerAddress() string {
	return "127.0.0.1:8999"
}

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

func createHTTP3Server(addr string, handler http.Handler) http3.Server {
	return http3.Server{
		Addr:    addr,
		Handler: handler,
	}
}

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
