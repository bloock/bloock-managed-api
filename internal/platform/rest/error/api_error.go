package api_error

import "net/http"

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (a APIError) Error() string {
	return a.Message
}

func NewAPIError(status int, message string) *APIError {
	return &APIError{Status: status, Message: message}
}
func NewBadRequestAPIError(message string) *APIError {
	return &APIError{Status: http.StatusBadRequest, Message: message}
}
func NewUnauthorizedAPIError(message string) *APIError {
	return &APIError{Status: http.StatusUnauthorized, Message: message}
}
func NewInternalServerAPIError(message string) *APIError {
	return &APIError{Status: http.StatusInternalServerError, Message: message}
}
