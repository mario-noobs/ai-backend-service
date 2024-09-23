package models

type BasicResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Response codes
const (
	Success            = "0000" // OK
	InvalidParamsErr   = "201"  // Created
	NullErr            = "204"  // No Content
	BadRequest         = "400"  // Bad Request
	Unauthorized       = "401"  // Unauthorized
	Forbidden          = "403"  // Forbidden
	JsonParseErr       = "404"  // Not Found
	NetworkErr         = "5000" // Internal Server Error
	ServiceUnavailable = "503"  // Service Unavailable
	Unknown            = "999"  // Service Unavailable
)

// Error messages
var errorMessages = map[string]string{
	Success:          "Success",
	InvalidParamsErr: "Access is denied due to invalid credentials.",
	NullErr:          "The requested resource could not be found.",
	BadRequest:       "A database error occurred.",
	Unknown:          "An unknown error occurred.",
}

// SetErrorMessage returns the error message corresponding to the error code
func SetErrorMessage(code string) BasicResponse {
	if msg, exists := errorMessages[code]; exists {
		return BasicResponse{
			Code:    code,
			Message: msg,
		}
	}
	return BasicResponse{
		Code:    Unknown,
		Message: errorMessages[Unknown],
	}
}

func SetErrorCodeMessage(code string, message string) BasicResponse {
	return BasicResponse{
		Code:    code,
		Message: message,
	}
}
