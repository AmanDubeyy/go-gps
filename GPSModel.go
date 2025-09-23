
package main

type GPSData struct {
	VehicleID string  `json:"vehicle_id"`
	RouteID   string  `json:"route_id"`
	City      string  `json:"city"`
	Status    string  `json:"status"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Speed     float64 `json:"speed"`
	Heading   int     `json:"heading"`
	DriverID  string  `json:"driver_id"`
}