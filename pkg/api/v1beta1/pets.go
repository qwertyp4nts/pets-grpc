package v1beta1

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"time"

	"github.com/qwertyp4nts/pets-grpc/pkg/integration/restapiprovider"
	proto "github.com/qwertyp4nts/pets-grpc/proto/v1beta1/pets"
)

// GetPet ...
func (s *Service) GetPet(ctx context.Context, req *proto.GetPetRequest) (*proto.GetPetResponse, error) {

	pets, err := s.adapters.RESTAPIProvider.GetPet(mapGetPetRequest(req))
	if err != nil {
		return nil, err
	}

	res, err := mapToPetResponse(ctx, pets)
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
		LastUpdated: timestamppb.New(time.Now().UTC()),
		Pet:         res,
	}, nil
}

func mapGetPetRequest(req *proto.GetPetRequest) restapiprovider.GetPetRequest {
	return restapiprovider.GetPetRequest{
		Type:  req.GetType(),
		Breed: req.GetBreed(),
	}
}

// GetPet ...
func (s *Service) GetPets(ctx context.Context, req *proto.GetPetsRequest) (*proto.GetPetsResponse, error) {

	pets, err := s.adapters.RESTAPIProvider.GetPets()
	if err != nil {
		return nil, err
	}

	p := []*proto.Pet{}

	for _, pet := range pets {
		p = append(p, &proto.Pet{
			Id:    pet.Id,
			Type:  pet.Type,
			Breed: pet.Breed,
			Risk:  pet.Risk,
		})
	}

	return &proto.GetPetsResponse{
		Pet: p,
	}, nil
}

func mapToPetResponse(ctx context.Context, pets *restapiprovider.Pets) (*proto.Pet, error) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println(pets.Id)
	log.Println(pets.Type)
	log.Println(pets.Breed)
	log.Println(pets.Risk)
	return &proto.Pet{
		Id:    pets.Id,
		Type:  pets.Type,
		Breed: pets.Breed,
		Risk:  pets.Risk,
	}, nil
}

// AddPet ...
func (s *Service) AddPet(ctx context.Context, req *proto.AddPetRequest) (*proto.AddPetResponse, error) {
	downstreamReq := restapiprovider.AddPetRequest{
		Breed: req.Breed,
		Type:  req.Type,
		Risk:  req.Risk,
	}
	err := s.adapters.RESTAPIProvider.AddPet(downstreamReq)
	if err != nil {
		return nil, err
	}

	return &proto.AddPetResponse{}, nil
}
