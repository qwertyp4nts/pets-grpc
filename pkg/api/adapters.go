package api

import (
	"context"
	"github.com/qwertyp4nts/pets-grpc/cmd/pets/config"
	"github.com/qwertyp4nts/pets-grpc/pkg/integration/restapiprovider"
)

// Adapters defines dependencies for the Pets API to use.
type Adapters struct {
	RESTAPIProvider restapiprovider.Servicer
}

// NewAdapters constructs a new set of dependencies to be used by the Pets API.
func NewAdapters(ctx context.Context, cfg *config.Config) (*Adapters, error) {
	restApiService := &restapiprovider.Service{
		AppSpec: cfg.AppSpec,
	}
	return &Adapters{
		RESTAPIProvider: restApiService,
	}, nil
}
