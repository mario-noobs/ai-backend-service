package service

import (
	"encoding/json"
	"errors"
	helper "golang-ai-management/helpers"
	"golang-ai-management/models"
	"golang-ai-management/models/response"
	"golang-ai-management/utils"
	"time"
)

type FaceService interface {
	ListIdentities() response.FaceResponse
	Enroll(face models.Face) response.FaceRegResponse
	Recognize(face models.Face) response.FaceRegResponse
}

type MarioFaceService struct {
	config MarioFaceServiceConfig
}

type userNames struct {
	Names []string `json:"names"`
}

// ListIdentities returns a list of identities.
func (s *MarioFaceService) ListIdentities() response.FaceResponse {
	var config = s.config.LoadMarioFaceServiceConfig()
	var code = models.BasicResponse{}
	var data = []string{}
	resp, err := helper.GetAPI(config.Host+config.listPath, make(map[string]string))
	if err != nil {
		code = models.SetErrorCodeMessage(models.NetworkErr, err.Error())
	}

	var responseJson userNames
	if err := json.Unmarshal(resp, &responseJson); err != nil {
		// Handle the error when unmarshalling JSON fails
		code = models.SetErrorCodeMessage(models.NetworkErr, err.Error())
	}

	// Check if the list is empty
	if len(responseJson.Names) == 0 {
		code = models.SetErrorCodeMessage(models.Success, errors.New("no identities found").Error())
	} else {
		code = models.SetErrorMessage(models.Success)
		data = responseJson.Names
	}

	return response.FaceResponse{
		code,
		data,
	}
}

// enroll registers a new face.
func (s *MarioFaceService) Enroll(face models.Face) response.FaceRegResponse {
	config := s.config.LoadMarioFaceServiceConfig()
	var code = models.BasicResponse{}

	payload, err := utils.StructToMap(face)
	if err != nil {
		code = models.SetErrorCodeMessage(models.InvalidParamsErr, err.Error())
	}

	resp, err := helper.PostAPI(config.Host+config.enrollPath, payload)
	if err != nil {
		code = models.SetErrorCodeMessage(models.NetworkErr, err.Error())
	}

	var result = models.BasicResponse{}
	if err := json.Unmarshal(resp, &result); err != nil {
		// Handle the error when unmarshalling JSON fails
		code = models.SetErrorCodeMessage(models.InvalidParamsErr, err.Error())
	}

	code = models.SetErrorCodeMessage(result.Code, result.Message)

	// Check if the list is empty
	if (result.Code != models.Success) && (result.Code != "") {
		code = models.SetErrorCodeMessage(result.Code, result.Message)
	}

	return response.FaceRegResponse{code, response.FaceData{
		Name:      face.Name,
		CreatedAt: time.Now().Format(time.RFC3339),
	}}
}

// recognize identifies a face.
func (s *MarioFaceService) Recognize(face models.Face) response.FaceRegResponse {

	// For example, return a mock response
	return response.FaceRegResponse{}
}
