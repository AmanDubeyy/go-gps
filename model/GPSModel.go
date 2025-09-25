package model

type GPSData struct {
	VehicleID   string  `json:"vehicle_id"`
	VehicleType string  `json:"vehicle_type"`
	RouteID     string  `json:"route_id"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Speed       float64 `json:"speed"`
	RegNo       string  `json:"reg_no"`
}
