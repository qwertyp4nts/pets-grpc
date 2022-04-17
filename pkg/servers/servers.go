package servers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCRegistration is the signature for gRPC registrations against a server.
type GRPCRegistration func(server *grpc.Server) error

// GRPCServer starts a gRPC server.
func GRPCServer(
	ctx context.Context,
	name string,
	host string,
	port uint16,
	registrations []GRPCRegistration,
) error {
	serverErr := make(chan error)
	defer close(serverErr)

	server := grpc.NewServer()
	reflection.Register(server)

	for _, register := range registrations {
		if err := register(server); err != nil {
			return errors.New("failed to register gRPC service")
		}
	}

	serverAddress := fmt.Sprintf("%s:%d", host, port)

	log.Printf("Starting %s gRPC server on %s", name, serverAddress)

	go startServer(server, serverAddress, serverErr)

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Printf("gracefully shutting down %s gRPC server", name)
		server.GracefulStop()
	}

	return <-serverErr
}

// HTTPServer starts a http server.
func HTTPServer(ctx context.Context, mux http.Handler, name string, host string, port uint16) error {
	serverErr := make(chan error)
	defer close(serverErr)

	serverAddress := fmt.Sprintf("%s:%d", host, port)
	server := &http.Server{Addr: serverAddress, Handler: mux}

	log.Printf("Starting %s http server on %s", name, serverAddress)

	// Emulate the behaviour of the gRPC server which will not error when shutdown is called twice
	go func(server *http.Server, serverErr chan error) {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			serverErr <- errors.New("could not start http server")
		} else {
			serverErr <- nil
		}
	}(server, serverErr)

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Printf("gracefully shutting down %s http server", name)
		_ = server.Shutdown(ctx)
	}

	return <-serverErr
}

func startServer(server *grpc.Server, serverAddress string, serverErr chan error) {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		serverErr <- errors.New("opening of TCP port failed")

		return
	}

	err = server.Serve(listener)
	if err != nil {
		serverErr <- errors.New("could not start listener for gRPC server")
	} else {
		serverErr <- nil
	}
}
