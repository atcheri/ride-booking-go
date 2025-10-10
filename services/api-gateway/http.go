package main

import (
	"encoding/json"
	"net/http"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var body tripPreviewRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// validation
	if body.UserId == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, "ok")

}
