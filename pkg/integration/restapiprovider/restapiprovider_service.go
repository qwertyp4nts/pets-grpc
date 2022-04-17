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
	"os"
)

// Servicer provides the transport-agnostic API for RestApiProvider.
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
	// Client  *http.Client // - For control over HTTP client headers, redirect policy, and other settings, create a Client
}

// GetPet ...
//func (s *Service) GetPets(ctx context.Context) (Pets, error) {
//	//payload := getPetsRequest()
//
//	//requestBody, err := json.Marshal(payload)
//	//if err != nil {
//	//	return fmt.Errorf("an error occurred calling the restapiprovider API - Error: '%v', Status: '%d'", err, 500)
//	//}
//
//	return s.getPets(), nil
//}

// getPets executes a call to restapiprovider's Get Pets API.
func (s *Service) GetPet(req GetPetRequest) (*Pets, error) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

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
		return nil, errors.New("error performing GET request to restapiprovider API")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Ignore errors as this is just for logging purposes
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("response from restapiprovider =/= OK. respBody= %v", body)
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

// getPets executes a call to restapiprovider's Get Pets API.
func (s *Service) GetPets() ([]*Pets, error) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	//path := fmt.Sprintf("%s/pets/dog/dalmatian", s.AppSpec.RestAPIProvider.Host)
	path := fmt.Sprintf("%s/pets", s.AppSpec.RestAPIProvider.Host)
	log.Println(path)

	resp, err := http.Get(path)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error performing GET request to restapiprovider API")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Ignore errors as this is just for logging purposes
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("response from restapiprovider =/= OK. respBody= %v", body)
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

	//objSlice, ok := createGetPetsResponse[0].([]interface{})

	return createGetPetsResponse, nil
}

// AddPet ...
func (s *Service) AddPet(req AddPetRequest) error {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	if req.Type == "" {
		log.Println("error - req.Type is required")
		return errors.New("error - req.Type is required")
	}

	if req.Breed == "" {
		log.Println("error - req.Breed is required")
		return errors.New("error - req.Breed is required")
	}

	//if req.Risk == 0 {
	//	log.Println("error - req.Risk is required")
	//	return errors.New("error - req.Risk is required")
	//} // cant tell if its 0 or not provided

	path := fmt.Sprintf("%s/pets/add", s.AppSpec.RestAPIProvider.Host)
	log.Println(path)

	requestBody, err := json.Marshal(req)

	resp, err := http.Post(path, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
		return errors.New("error performing GET request to restapiprovider API")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// Ignore errors as this is just for logging purposes
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("response from restapiprovider =/= OK. respBody= %v", body)
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
