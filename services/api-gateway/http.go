package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/atcheri/ride-booking-go/shared/contracts"
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

	jsonBody, _ := json.Marshal(body)
	reader := bytes.NewReader(jsonBody)
	resp, err := http.Post("http://trip-service:8083/preview", "application-json", reader)
	if err != nil {
		writeJSON(w, http.StatusNotFound, contracts.APIError{
			Code:    strconv.Itoa(http.StatusNotFound),
			Message: err.Error(),
		})
		return
	}

	defer resp.Body.Close()

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "failed to parse JSON data from trip service", http.StatusBadRequest)
		return
	}

	response := contracts.APIResponse{Data: respBody}

	writeJSON(w, http.StatusCreated, response)
}
