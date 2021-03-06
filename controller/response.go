package controller

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message,omitempty"`
}

type CreatedResponse struct {
	Message string `json:"message,omitempty"`
	Id      string `json:"id,omitempty"`
}

type ErrResp struct {
	Message string `json:"error,omitempty"`
}

type JWTResponse struct {
	AccessToken string `json:"access_token,omitempty"`
}

type ValidationError struct {
	Message string
	Status  int
}

func (e *ValidationError) Error() string {
	return e.Message
}

// JSONResp - respond json with status
func JSONResp(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
