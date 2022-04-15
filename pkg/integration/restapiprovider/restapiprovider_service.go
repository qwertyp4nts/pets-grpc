package restapiprovider

import (
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
	GetPets() (*Pets, error)
}

// Pets holds the pet attributes used in a GetPetResponse.
type Pets struct {
	Type string `json:"type"`
	//Addresses    []GetPartyResponseAddress `json:"addresses"`
	//DateOfBirth  string                    `json:"dateOfBirth"`
	//Emails       []Email                   `json:"emails"`
	//Phones       []Phone                   `json:"phones"`
	//Identifiers  []IdentifierInfo          `json:"identifiers"`
	//TaxCountries []TaxCountry              `json:"taxCountries"`
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
func (s *Service) GetPets() (*Pets, error) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	path := fmt.Sprintf("%s/pets/dog/dalmatian", s.AppSpec.RestAPIProvider.Host)
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
