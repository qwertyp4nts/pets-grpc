package startup

import (
	apiV1beta1 "github.com/qwertyp4nts/pets-grpc/pkg/api/v1beta1"
	"github.com/qwertyp4nts/pets-grpc/pkg/servers"
	protov1beta1 "github.com/qwertyp4nts/pets-grpc/proto/v1beta1/pets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

// GRPCRegistrations returns functions with implementations of registrations against the gRPC server.
// nolint:funlen
func GRPCRegistrations(healthServer *health.Server, apiV1beta1 apiV1beta1.Servicer) []servers.GRPCRegistration {
	registrations := []servers.GRPCRegistration{
		// Fabric Self Service v1beta1 registration
		func(server *grpc.Server) error {
			protov1beta1.RegisterPetsAPIServer(server, apiV1beta1)
			return nil
		},
		// Health Check
		func(server *grpc.Server) error {
			healthgrpc.RegisterHealthServer(server, healthServer)
			return nil
		},
	}

	return registrations
}
