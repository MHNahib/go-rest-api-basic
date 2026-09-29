package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ErrorResponseType struct {
	Type string `json:"error_type"`
}
type ErrorResponse struct {
	Error ErrorResponseType `json:"error"`
}

func NewResponse(status int, message string, data interface{}) Response {
	return Response{
		Status:  !(status >= http.StatusBadRequest),
		Message: message,
		Data:    data,
	}
}

func WriteJson(w http.ResponseWriter, statusCode int, data interface{}, message string) error {
	w.Header().Set("Content-Type", "application/json")

	_message := message
	if message == "" {
		_message = http.StatusText(statusCode)
	}

	response := NewResponse(statusCode, _message, data)

	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(response)
}

func GenericError(err error) ErrorResponse {
	return ErrorResponse{Error: ErrorResponseType{Type: err.Error()}}
}

func ValidationErrors(errs validator.ValidationErrors) ErrorResponse {
	var errMessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, strings.ToLower(err.Field())+" is required")

		default:
			errMessages = append(errMessages, strings.ToLower(err.Field())+" is invalid")
		}
	}

	return ErrorResponse{Error: ErrorResponseType{Type: strings.Join(errMessages, ", ")}}

}
