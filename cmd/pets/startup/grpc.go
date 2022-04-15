package startup

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	apiV1beta1 "github.com/qwertyp4nts/pets-grpc/pkg/api/v1beta1"
	"github.com/qwertyp4nts/pets-grpc/pkg/servers"
	protov1beta1 "github.com/qwertyp4nts/pets-grpc/proto/v1beta1/pets"
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
			// healthServer.GRPC.RegisterWith(server)
			// healthServer.SetReady(true)
			return nil
		},
	}

	return registrations
}
