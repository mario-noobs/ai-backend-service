package response

import "golang-ai-management/models"

type FaceData struct {
	Name        *string  `json:"name"`
	Probability *float64 `json:"probability"`
	CreatedAt   string   `json:"created_at"`
}

type FaceRegResponse struct {
	models.BasicResponse
	Data FaceData `json:"data"`
}

type FaceResponse struct {
	models.BasicResponse
	Data []string `json:"data"`
}
