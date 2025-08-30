package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func GetSuccessResponseWithBody(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(body)
	if err != nil {
		log.Printf("http.GetSuccessResponseWithBody: %v\n", err)
	}
}

func GetErrorResponse(w http.ResponseWriter, status int, err error, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errEnc := json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err.Error(),
		Details: details,
	})
	if errEnc != nil {
		log.Printf("http.GetErrorResponse: %v\n", err)
	}
}
