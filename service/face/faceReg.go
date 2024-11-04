package face

import (
	"context"
	"encoding/json"
	helper "golang-ai-management/helpers"
	"golang-ai-management/models"
	"golang-ai-management/models/response"
	"golang-ai-management/utils"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func NewFaceBusiness(faceService FaceService, config MarioFaceServiceConfig, time helper.Timer) *FaceBussiness {
	return &FaceBussiness{
		faceService,
		config,
		time,
	}
}

type FaceService interface {
	Enroll(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse
	Recognize(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse
}

type FaceBussiness struct {
	FaceBussiness FaceService
	config        MarioFaceServiceConfig
	time          helper.Timer
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func (f FaceBussiness) Enroll(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse {

	var method = "FaceBussiness_Enroll"
	f.time.Start()
	logger.Info("request", "requestId", face.TransactionId, "method", method)

	cfg := f.config.LoadMarioFaceServiceConfig()

	payload, err := utils.StructToMap(face)
	if err != nil {
		logger.Error("response", "method", method, "requestId", face.TransactionId, "err", err, "ms", nil)

	}
	logger.DebugContext(ctx, method, "URLs", "requestId", face.TransactionId, cfg.Host+cfg.listPath)

	resp, err := helper.PostAPI(cfg.Host+cfg.enrollPath, payload, jwt)
	if err != nil {
		logger.Error("response", "method", method, "requestId", face.TransactionId, "err", err, "ms", nil)

	}

	result, err := MapResponse(resp)
	if err != nil {
		// Handle the error when unmarshalling JSON fails
		logger.Error("response", "method", method, "requestId", face.TransactionId, "err", err, "ms", nil)

	}

	logger.Info("response", "method", method, "requestId", face.TransactionId, "data", result, "ms", f.time.End())

	return response.FaceRegResponse{UserId: *face.Name, BasicResponse: result.BasicResponse, Data: response.FaceData{
		Name:      face.Name,
		CreatedAt: time.Now().Format(time.RFC3339),
	}}
}

func (f FaceBussiness) Recognize(ctx context.Context, face models.Face, jwt string) response.FaceRegResponse {
	var method = "FaceBussiness_Enroll"
	f.time.Start()
	logger.Info("request", "method", method, "requestId", face.TransactionId)

	cfg := f.config.LoadMarioFaceServiceConfig()
	var code = models.BasicResponse{}

	payload, err := utils.StructToMap(face)
	if err != nil {
		logger.Error("response", "method", method, "requestId", face.TransactionId, "err", err, "ms", nil)
		code = models.SetErrorCodeMessage(models.InvalidParamsErr, err.Error())
	}

	resp, err := helper.PostAPI(cfg.Host+cfg.recognizePath, payload, jwt)
	if err != nil {
		logger.Error("response", "method", method, "requestId", face.TransactionId, "err", err, "ms", nil)
		code = models.SetErrorCodeMessage(models.NetworkErr, err.Error())
	}

	result, err := MapResponse(resp)
	if err != nil {
		// Handle the error when unmarshalling JSON fails
		logger.Error("response", "method", method, "requestId", face.TransactionId, "err", err, "ms", nil)
		code = models.SetErrorCodeMessage(models.BadRequest, err.Error())
	}

	if result.Code == models.Success {
		logger.Info("response", "method", method, "requestId", face.TransactionId, "data", result, "ms", f.time.End())
		return result
	} else {
		logger.Info("response", "method", method, "requestId", face.TransactionId, "data", code, "ms", f.time.End())
		return response.FaceRegResponse{BasicResponse: code, Data: response.FaceData{
			CreatedAt: time.Now().Format(time.RFC3339),
		}}
	}
}

func MapResponse(jsonData []byte) (response.FaceRegResponse, error) {

	type SearchData struct {
		Confidence map[string]string `json:"searh_result"`
		UserId     string            `json:"userId"`
	}

	var tempResponse struct {
		Code      string     `json:"code"`
		Message   string     `json:"message"`
		RequestId string     `json:"requestId"`
		UserId    string     `json:"userId"`
		RawImage  string     `json:"rawImage"`
		Flow      string     `json:"flow"`
		Data      SearchData `json:"searchData"`
	}

	// Unmarshal the JSON response into the struct
	err := json.Unmarshal(jsonData, &tempResponse)
	if err != nil {
		return response.FaceRegResponse{}, err
	}

	// Prepare the FaceRegResponse
	regResponse := response.FaceRegResponse{
		UserId: tempResponse.UserId,
		BasicResponse: models.BasicResponse{
			Code:      tempResponse.Code,
			Message:   tempResponse.Message,
			RequestId: tempResponse.RequestId,
		},
		Data: response.FaceData{
			Image:     &tempResponse.RawImage,
			CreatedAt: time.Now().Format(time.RFC3339), // Example created_at date; set as needed
		},
	}

	// Handle confidence mapping for any name
	for name, probStr := range tempResponse.Data.Confidence {
		probability, err := strconv.ParseFloat(string(probStr), 64)
		if err == nil {
			regResponse.Data.Name = &name
			regResponse.Data.Probability = &probability
			break // Remove this if you want to capture all confidence entries
		}
	}

	return regResponse, nil
}
