package model

type VehicleRoute struct {
	VehicleID string  `json:"vehicle_id"`
	City      string  `json:"city"`
	Status    string  `json:"status"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	DriverID  string  `json:"driver_id"`
	Time      string  `json:"time"`
}
