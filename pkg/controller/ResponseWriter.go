package controller

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, status int, redirect string, flashMessage string, flashType string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"redirect":     redirect,
		"flashMessage": flashMessage,
		"flashType":    flashType,
	})
}