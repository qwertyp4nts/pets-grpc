package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/health"

	"github.com/qwertyp4nts/pets-grpc/cmd/pets/config"
	"github.com/qwertyp4nts/pets-grpc/cmd/pets/config/app"
	// "github.com/qwertyp4nts/pets-grpc/cmd/pets/metadata"
	"github.com/qwertyp4nts/pets-grpc/cmd/pets/startup"
	"github.com/qwertyp4nts/pets-grpc/pkg/api"
	"github.com/qwertyp4nts/pets-grpc/pkg/api/v1beta1"
	"github.com/qwertyp4nts/pets-grpc/pkg/servers"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(ctx)
	if err != nil {
		fmt.Errorf("failed to load app config. error: %v", err)
	}

	// Create health server
	healthServer := health.NewServer()

	errChan := make(chan error)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	adapters, err := api.NewAdapters(ctx, cfg)
	if err != nil {
		fmt.Errorf("failed to create adapters. error: %v", err)
	}

	// Run main gRPC server
	go func() {
		errChan <- runGRPCServer(ctx, &cfg.AppSpec, healthServer, adapters)
	}()

	// Run operations server
	go func() {
		errChan <- runOperationsServer(ctx, cfg, healthServer)
	}()

	select {
	case err = <-errChan:
		fmt.Errorf("server returned error. error: %v", err)
	case <-signals:
		fmt.Errorf("%s terminated by SIGTERM", cfg.AppSpec.AppName)
	}

	// clean up before quit
	close(errChan)
	close(signals)
}

func runGRPCServer(
	ctx context.Context,
	appSpec *app.Spec,
	healthServer *health.Server,
	adapters *api.Adapters,
) error {
	petsv1beta1 := v1beta1.NewService(adapters, appSpec)
	registrations := startup.GRPCRegistrations(healthServer, petsv1beta1)

	return servers.GRPCServer(ctx, appSpec.AppName, appSpec.Host, appSpec.Port, registrations)
}

func runOperationsServer(ctx context.Context, cfg *config.Config, healthServer *health.Server) error {
	serverName := fmt.Sprintf("%s operations", cfg.AppSpec.AppName)

	mux := http.NewServeMux()

	// healthServer.HTTP.RegisterWith(mux)

	//if err := otelhealth.Register(healthServer.State); err != nil {
	//	log.Error(ctx, err).Msgf("Failed to register health metrics")
	//
	//	return anzErrors.New(
	//		codes.Internal,
	//		fmt.Sprintf("%s server could not register health metrics: %v", serverName, err),
	//		anzErrors.NewErrorInfo(ctx, errcodes.Unknown, "failed to register health metrics"),
	//	)
	//}

	return servers.HTTPServer(ctx, mux, serverName, cfg.OpsSpec.Host, cfg.OpsSpec.Port)
}
