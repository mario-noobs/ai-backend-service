package face

import (
	"context"
	"encoding/json"
	"golang-ai-management/config"
	helper "golang-ai-management/helpers"
	"golang-ai-management/logger"
	"golang-ai-management/models"
	"golang-ai-management/models/response"
	"golang-ai-management/utils"
	"strconv"
	"time"
)

func NewFaceBusiness(faceService FaceService, config MarioFaceServiceConfig) *FaceBussiness {
	return &FaceBussiness{
		faceService,
		config,
	}
}

type FaceService interface {
	Enroll(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse
	Recognize(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse
}

type FaceBussiness struct {
	FaceBussiness FaceService
	config        MarioFaceServiceConfig
}

var serverConfig = config.Config
var factory = logger.LoggerFactory{}
var newLogger, err = factory.NewLogger(serverConfig.LogType, serverConfig.LogFormat, serverConfig.LogLevel)

func (f FaceBussiness) Enroll(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse {
	config := f.config.LoadMarioFaceServiceConfig()
	var code = models.BasicResponse{}
	newLogger.DebugArgs("Enroll", "Params", face)

	payload, err := utils.StructToMap(face)
	if err != nil {
		newLogger.Error(err.Error())
		code = models.SetErrorCodeMessage(models.InvalidParamsErr, err.Error())
	}
	newLogger.DebugArgs("Enroll", "URLs", config.Host+config.listPath)

	resp, err := helper.PostAPI(config.Host+config.enrollPath, payload, jwt)
	if err != nil {
		newLogger.Error(err.Error())
		code = models.SetErrorCodeMessage(models.NetworkErr, err.Error())
	}

	var result = models.BasicResponse{}
	if err := json.Unmarshal(resp, &result); err != nil {
		// Handle the error when unmarshalling JSON fails
		newLogger.Error(err.Error())
		code = models.SetErrorCodeMessage(models.InvalidParamsErr, err.Error())
	}

	// Check if the list is empty
	if (result.Code != models.Success) && (result.Code != "") {
		code = models.SetErrorCodeMessage(result.Code, result.Message)
	}

	newLogger.DebugArgs("Enroll", "response", code)

	return response.FaceRegResponse{code, response.FaceData{
		Name:      face.Name,
		CreatedAt: time.Now().Format(time.RFC3339),
	}}
}

func (f FaceBussiness) Recognize(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse {
	config := f.config.LoadMarioFaceServiceConfig()
	var code = models.BasicResponse{}
	newLogger.DebugArgs("Recognize", "Params", face)

	payload, err := utils.StructToMap(face)
	if err != nil {
		newLogger.Error(err.Error())
		code = models.SetErrorCodeMessage(models.InvalidParamsErr, err.Error())
	}

	resp, err := helper.PostAPI(config.Host+config.recognizePath, payload, jwt)
	if err != nil {
		newLogger.Error(err.Error())
		code = models.SetErrorCodeMessage(models.NetworkErr, err.Error())
	}

	result, err := MapResponse(resp)
	if err != nil {
		// Handle the error when unmarshalling JSON fails
		newLogger.Error(err.Error())
		code = models.SetErrorCodeMessage(models.BadRequest, err.Error())
	}

	if result.Code == models.Success {
		newLogger.DebugArgs("Recognize", "response", result.BasicResponse)
		return result
	} else {
		newLogger.DebugArgs("Recognize", "response", code)
		return response.FaceRegResponse{BasicResponse: code, Data: response.FaceData{
			CreatedAt: time.Now().Format(time.RFC3339),
		}}
	}
}

func MapResponse(jsonData []byte) (response.FaceRegResponse, error) {
	var tempResponse struct {
		Code       string            `json:"code"`
		Confidence map[string]string `json:"confidence"`
		RawImage   string            `json:"raw_image"`
	}

	// Unmarshal the JSON response into the struct
	err := json.Unmarshal(jsonData, &tempResponse)
	if err != nil {
		return response.FaceRegResponse{}, err
	}

	// Prepare the FaceRegResponse
	regResponse := response.FaceRegResponse{
		BasicResponse: models.SetErrorMessage(models.Success),
		Data: response.FaceData{
			Image:     &tempResponse.RawImage,
			CreatedAt: time.Now().Format(time.RFC3339), // Example created_at date; set as needed
		},
	}

	// Handle confidence mapping for any name
	for name, probStr := range tempResponse.Confidence {
		probability, err := strconv.ParseFloat(probStr, 64)
		if err == nil {
			regResponse.Data.Name = &name
			regResponse.Data.Probability = &probability
			break // Remove this if you want to capture all confidence entries
		}
	}

	return regResponse, nil
}
