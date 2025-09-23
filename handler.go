package main

import (
	"encoding/json"
	"net/http"

	"strconv"
	"time"
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

		var data GPSData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		select {
		case GpsChan <- data:
			w.WriteHeader(http.StatusAccepted)
		default:
			http.Error(w, "server overloaded", http.StatusServiceUnavailable)
		}
	}
}

func RouteHandler(repo *InfluxGPSRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vehicleID := r.URL.Query().Get("vehicle_id")
		start, _ := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
		end, _ := time.Parse(time.RFC3339, r.URL.Query().Get("end"))

		route, err := repo.GetRoute(vehicleID, start, end)
		if err != nil {
			http.Error(w, "failed to get route: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(route)
	}
}

func SearchHandler(repo *InfluxGPSRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minLat, _ := strconv.ParseFloat(r.URL.Query().Get("minLat"), 64)
		maxLat, _ := strconv.ParseFloat(r.URL.Query().Get("maxLat"), 64)
		minLon, _ := strconv.ParseFloat(r.URL.Query().Get("minLon"), 64)
		maxLon, _ := strconv.ParseFloat(r.URL.Query().Get("maxLon"), 64)

		points, err := repo.SearchVehicles(minLat, maxLat, minLon, maxLon)
		if err != nil {
			http.Error(w, "failed to search vehicles: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(points)
	}
}
