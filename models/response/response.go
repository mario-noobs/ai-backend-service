package response

import "golang-ai-management/models"

type FaceData struct {
	Name        *string  `json:"name"`
	Probability *float64 `json:"probability"`
	CreatedAt   string   `json:"created_at"`
	Image       *string  `json:"image"`
}

type FaceRegResponse struct {
	UserId string `json:"userId"`
	models.BasicResponse
	Data FaceData `json:"data"`
}

type FaceResponse struct {
	models.BasicResponse
	Data []string `json:"data"`
}
