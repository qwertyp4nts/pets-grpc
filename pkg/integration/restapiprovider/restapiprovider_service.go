package restapiprovider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qwertyp4nts/pets-grpc/cmd/pets/config/app"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Servicer provides the transport-agnostic API for RESTAPIProvider.
type Servicer interface {
	GetPet(GetPetRequest) (*Pets, error)
	GetPets() ([]*Pets, error)
	AddPet(AddPetRequest) error
}

type GetPetRequest struct {
	Type  string `json:"type"`
	Breed string `json:"breed"`
}

// Pets holds the pet attributes used in a GetPetResponse.
type Pets struct {
	Id    string `json:"_id"`
	Type  string `json:"type"`
	Breed string `json:"breed"`
	Risk  int32  `json:"risk"`
}

type AddPetRequest struct {
	Type  string `json:"type"`
	Breed string `json:"breed"`
	Risk  int32  `json:"risk"`
}

// Service holds the RESTAPIProvider Service attributes.
type Service struct {
	AppSpec app.Spec
	// Client  *http.Client // - For control over HTTP client headers, redirect policy, and other settings, create a Client here
}

// GetPet executes a call to RESTAPIProvider's Get Pets API.
func (s *Service) GetPet(req GetPetRequest) (*Pets, error) {
	if req.Type == "" {
		log.Println("error - req.Type is required")
		return nil, errors.New("error - req.Type is required")
	}

	if req.Breed == "" {
		log.Println("error - req.Breed is required")
		return nil, errors.New("error - req.Breed is required")
	}

	path := fmt.Sprintf("%s/pets/%s/%s", s.AppSpec.RestAPIProvider.Host, req.Type, req.Breed)
	log.Println(path)

	resp, err := http.Get(path)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error performing GET request to RESTAPIProvider API")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Ignore errors as this is just for logging purposes
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("response from RESTAPIProvider =/= OK. respBody= %v", body)
	}

	responseBody, responseErr := s.readResponse(resp.StatusCode, resp.Body)
	if responseErr != nil {
		return nil, responseErr
	}

	var createGetPetsResponse []*Pets

	err = json.Unmarshal(responseBody, &createGetPetsResponse)
	if err != nil {
		return nil, errors.New("server error")
	}

	return createGetPetsResponse[0], nil
}

// GetPets executes a call to RESTAPIProvider's Get Pets API.
func (s *Service) GetPets() ([]*Pets, error) {
	path := fmt.Sprintf("%s/pets", s.AppSpec.RestAPIProvider.Host)
	log.Println(path)

	resp, err := http.Get(path)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error performing GET request to RESTAPIProvider API")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Ignore errors as this is just for logging purposes
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("response from RESTAPIProvider =/= OK. respBody= %v", body)
	}

	responseBody, responseErr := s.readResponse(resp.StatusCode, resp.Body)
	if responseErr != nil {
		return nil, responseErr
	}

	var createGetPetsResponse []*Pets

	err = json.Unmarshal(responseBody, &createGetPetsResponse)
	if err != nil {
		return nil, errors.New("server error")
	}

	return createGetPetsResponse, nil
}

// AddPet ...
func (s *Service) AddPet(req AddPetRequest) error {
	if req.Type == "" {
		log.Println("error - req.Type is required")
		return errors.New("error - req.Type is required")
	}

	if req.Breed == "" {
		log.Println("error - req.Breed is required")
		return errors.New("error - req.Breed is required")
	}

	path := fmt.Sprintf("%s/pets/add", s.AppSpec.RestAPIProvider.Host)
	log.Println(path)

	requestBody, err := json.Marshal(req)

	resp, err := http.Post(path, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
		return errors.New("error performing GET request to RESTAPIProvider API")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// Ignore errors as this is just for logging purposes
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("response from RESTAPIProvider =/= OK. respBody= %v", body)
	}

	_, responseErr := s.readResponse(resp.StatusCode, resp.Body)
	if responseErr != nil {
		return responseErr
	}

	return nil
}

func (s *Service) readResponse(respSc int, respBody io.Reader) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, errors.New("server error")
	}

	if respSc < 200 || respSc > 299 {
		log.Printf("Error Response Body: '%s'", string(responseBody))

		return nil, errors.New("response from HTTP downstream NOT OK")
	}

	return responseBody, nil
}
