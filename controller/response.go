package controller

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message,omitempty"`
}

type ErrResp struct {
	Message string `json:"error,omitempty"`
}

type JWTResponse struct {
	AccessToken string `json:"access_token,omitempty"`
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
