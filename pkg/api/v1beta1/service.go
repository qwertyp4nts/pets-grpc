package v1beta1

import (
	proto "github.com/qwertyp4nts/pets-grpc/proto/v1beta1/pets"

	"github.com/qwertyp4nts/pets-grpc/cmd/pets/config/app"
	"github.com/qwertyp4nts/pets-grpc/pkg/api"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -build_flags=-mod=mod -destination=mocks/service.go -package=pets github.com/qwertyp4nts/pets-grpc/pkg/api/v1beta1 Servicer

// Servicer provides the transport-agnostic API for Pets
type Servicer interface {
	proto.PetsAPIServer
}

// Service holds onto the state of the Pets service.
type Service struct {
	adapters *api.Adapters
	appSpec  app.Spec

	proto.UnimplementedPetsAPIServer
}

// NewService creates an instance of a Pets Service.
func NewService(adapters *api.Adapters, appSpec *app.Spec) *Service {
	return &Service{
		adapters: adapters,
		appSpec:  *appSpec,
	}
}
