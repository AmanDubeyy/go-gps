package main

import (
	"encoding/json"
	"net/http"

	"ondc/model"
)

type GPSHandler struct {
	repo *InfluxGPSRepository
}

func NewGPSHandler(repo *InfluxGPSRepository) *GPSHandler {
	return &GPSHandler{repo: repo}
}

func InsertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var data model.GPSData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		select {
		case GpsChan <- data:
			json.NewEncoder(w).Encode(map[string]string{
				"code":    "success",
				"message": "GPS location updated",
			})
		default:
			http.Error(w, "server overloaded", http.StatusServiceUnavailable)
		}
	}
}

func SearchByID(repo *InfluxGPSRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vehicleID := r.URL.Query().Get("vehicle_id")

		route, err := repo.GetRoute(vehicleID)
		if err != nil {
			http.Error(w, "failed to get route: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(successResponse(route))
	}
}

func successResponse(data any) map[string]any{
	return map[string]any{
		"code" : "sucess",
		"data" : data,
	}
}