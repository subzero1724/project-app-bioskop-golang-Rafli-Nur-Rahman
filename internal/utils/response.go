package utils

import (
	"encoding/json"
	"net/http"

	"cinema-booking-system/internal/dto"
)

// RespondWithJSON writes a JSON response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// RespondWithSuccess sends a successful response
func RespondWithSuccess(w http.ResponseWriter, code int, data interface{}, message string) {
	response := dto.Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	RespondWithJSON(w, code, response)
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := dto.ErrorResponse{
		Success: false,
		Error:   message,
	}
	RespondWithJSON(w, code, response)
}

// RespondWithValidationError sends a validation error response
func RespondWithValidationError(w http.ResponseWriter, err error) {
	response := dto.ErrorResponse{
		Success: false,
		Error:   "Validation failed",
		Message: err.Error(),
	}
	RespondWithJSON(w, http.StatusBadRequest, response)
}
