package main

import (
	"context"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"ondc/model"
)

type InfluxGPSRepository struct {
	client *influxdb3.Client
}

func NewInfluxGPSRepository(client *influxdb3.Client) *InfluxGPSRepository {
	return &InfluxGPSRepository{client: client}
}

func (r *InfluxGPSRepository) InsertGPS(ctx context.Context, data model.GPSData) error {
	point := influxdb3.NewPoint(
		VEHICLE_LOCATION,
		map[string]string{
			"vehicle_id":   data.VehicleID,
			"vehicle_type": data.VehicleType,
			"route_id":     data.RouteID,
			"city":         data.City,
			"reg_no":       data.RegNo,
		},
		map[string]any{
			"lat":   data.Lat,
			"lon":   data.Lon,
			"speed": data.Speed,
		},
		time.Now(),
	)

	return r.client.WritePoints(ctx, []*influxdb3.Point{point})
}

func (r *InfluxGPSRepository) GetRoute(vehicleID string) ([]model.VehicleRoute, error) {
	ctx := context.Background()

	iter, err := r.client.Query(ctx, getByRouteID(vehicleID))
	if err != nil {
		return nil, err
	}

	var vehicles []model.VehicleRoute
	for iter.Next() {
		row := iter.Value()

		vehicles = append(vehicles, model.VehicleRoute{
			VehicleID: row["vehicle_id"].(string),
			VehicleType: row["vehicle_type"].(string),
			RegNo: row["reg_no"].(string),
			City: row["city"].(string),
			Lat:  row["lat"].(float64),
			Lon:  row["lon"].(float64),
			Time: row["time"].(time.Time).Format(time.RFC3339),
		})
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return vehicles, nil
}
