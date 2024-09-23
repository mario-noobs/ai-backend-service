package service

import (
	"encoding/json"
	"errors"
	"fmt"
	helper "golang-ai-management/helpers"
	"golang-ai-management/models"
	"golang-ai-management/models/basicModel"
)

type FaceService interface {
	ListIdentities() ([]string, error)
	enroll(face models.Face) (basicModel.Response, error)
	recognize(face models.Face) (basicModel.Response, error)
}

type MarioFaceService struct {
	config MarioFaceServiceConfig
}

type userNames struct {
	Names []string `json:"names"`
}

// ListIdentities returns a list of identities.
func (s *MarioFaceService) ListIdentities() ([]string, error) {
	var config = LoadMarioFaceServiceConfig()
	resp, err := helper.GetAPI(config.Host+config.listPath, make(map[string]string))
	if err != nil {
		// Handle the error when the request fails
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}

	// Parse the JSON response
	var response userNames
	if err := json.Unmarshal(resp, &response); err != nil {
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

	return basicModel.Response{Code: "true", Message: "Face recognized"}, nil
}

// recognize identifies a face.
func (s *MarioFaceService) recognize(face models.Face) (basicModel.Response, error) {

	// For example, return a mock response
	return basicModel.Response{Code: "true", Message: "Face recognized"}, nil
}
