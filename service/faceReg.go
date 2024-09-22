package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-ai-management/models"
	"golang-ai-management/models/basicModel"
	"io/ioutil"
	"net/http"
)

type FaceService interface {
	ListIdentities() ([]string, error)
	enroll(face models.Face) (basicModel.Response, error)
	recognize(face models.Face) (basicModel.Response, error)
}

type MarioFaceService struct {
}

type userNames struct {
	Names []string `json:"names"`
}

// ListIdentities returns a list of identities.
func (s *MarioFaceService) ListIdentities() ([]string, error) {
	// Define the API endpoint
	url := "http://75.119.149.223:5000/face/get-list"

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		// Handle the error when the request fails
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Handle the error when reading the body fails
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var response userNames
	if err := json.Unmarshal(body, &response); err != nil {
		// Handle the error when unmarshalling JSON fails
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Check if the list is empty
	if len(response.Names) == 0 {
		return nil, errors.New("no identities found")
	}

	return response.Names, nil
}

// enroll registers a new face.
func (s *MarioFaceService) enroll(face models.Face) (basicModel.Response, error) {

	// Add logic to save the face to your data store
	return basicModel.Response{Code: "true", Message: "Face recognized"}, nil
}

// recognize identifies a face.
func (s *MarioFaceService) recognize(face models.Face) (basicModel.Response, error) {

	// For example, return a mock response
	return basicModel.Response{Code: "true", Message: "Face recognized"}, nil
}
