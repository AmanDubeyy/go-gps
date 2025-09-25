package model

type VehicleRoute struct {
	VehicleID   int     `json:"vehicle_id"`
	VehicleType string  `json:"vehicle_type"`
	RegNo       string  `json:"reg_no"`
	City        string  `json:"city"`
	Status      string  `json:"status"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Time        string  `json:"time"`
}
