package model

type VehicleRoute struct {
	VehicleID   string  `json:"vehicle_id"`
	VehicleType string  `json:"vehicle_type"`
	RegNo       string  `json:"reg_no"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Time        string  `json:"time"`
}
