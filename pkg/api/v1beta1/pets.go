package v1beta1

import (
	"context"

	"github.com/qwertyp4nts/pets-grpc/pkg/integration/restapiprovider"
	proto "github.com/qwertyp4nts/pets-grpc/proto/v1beta1/pets"
)

// GetPet ...
func (s *Service) GetPet(ctx context.Context, req *proto.GetPetRequest) (*proto.GetPetResponse, error) {

	pets, err := s.adapters.RESTAPIProvider.GetPets()
	if err != nil {
		return nil, err
	}

	res, err := mapToPetsResponse(ctx, pets)
	if err != nil {
		return nil, err
		//return nil, anzErrors.New(
		//	codes.Internal,
		//	"Failed to map to GetPartyResponse",
		//	anzErrors.NewErrorInfo(ctx, errcodes.DownstreamFailure, anzErrors.GetMessage(err)),
		//	anzErrors.WithCause(err),
		//)
	}

	return &proto.GetPetResponse{
		Pet: res,
	}, nil
}

func mapToPetsResponse(ctx context.Context, pets *restapiprovider.Pets) (*proto.Pet, error) {
	return &proto.Pet{
		Type: pets.Type,
	}, nil
}

// UpdatePet ...
//func (s *Service) UpdatePet(ctx context.Context, req *proto.UpdatePetRequest) (*proto.UpdatePetResponse, error) {
//	return &proto.UpdatePartyResponse{}, nil
//}
