package main

import (
	"context"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

type InfluxGPSRepository struct {
	client *influxdb3.Client
	bucket string
}

func NewInfluxGPSRepository(client *influxdb3.Client, bucket string) *InfluxGPSRepository {
	return &InfluxGPSRepository{client: client, bucket: bucket}
}

func (r *InfluxGPSRepository) InsertGPS(ctx context.Context, data GPSData) error {
	point := influxdb3.NewPoint(
		VEHICLE_LOCATION,
		map[string]string{
			"vehicle_id": data.VehicleID,
			"route_id":   data.RouteID,
			"city":       data.City,
			"status":     data.Status,
			"driver_id":  data.DriverID,
		},
		map[string]interface{}{
			"lat":     data.Lat,
			"lon":     data.Lon,
			"speed":   data.Speed,
			"heading": data.Heading,
		},
		time.Now(),
	)

	return r.client.WritePoints(ctx, []*influxdb3.Point{point})
}

func (r *InfluxGPSRepository) GetRoute(vehicleID string, start, end time.Time) ([]GPSData, error) {
	ctx := context.Background()

	if start.IsZero() {
		start = time.Now().Add(-24 * time.Hour)
	}

	if end.IsZero() {
		end = time.Now()
	}

	iter, err := r.client.Query(ctx, getByRouteID(vehicleID))
	if err != nil {
		return nil, err
	}

	var vehicles []GPSData
	for iter.Next() {
		row := iter.Value()

		vehicles = append(vehicles, GPSData{
			VehicleID: row["vehicle_id"].(string),
			RouteID:   row["route_id"].(string),
			City:      row["city"].(string),
			Status:    row["status"].(string),
			Lat:       row["lat"].(float64),
			Lon:       row["lon"].(float64),
			Speed:     row["speed"].(float64),
			Heading:   int(row["heading"].(int64)),
			DriverID:  row["driver_id"].(string),
		})
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return vehicles, nil
}

func (r *InfluxGPSRepository) SearchVehicles(minLat, maxLat, minLon, maxLon float64) ([]GPSData, error) {
	ctx := context.Background()

	iter, err := r.client.Query(ctx, getAll)
	if err != nil {
		return nil, err
	}

	var vehicles []GPSData
	for iter.Next() {
		row := iter.Value()

		vehicles = append(vehicles, GPSData{
			VehicleID: row["vehicle_id"].(string),
			RouteID:   row["route_id"].(string),
			City:      row["city"].(string),
			Status:    row["status"].(string),
			Lat:       row["lat"].(float64),
			Lon:       row["lon"].(float64),
			Speed:     row["speed"].(float64),
			Heading:   int(row["heading"].(int64)),
			DriverID:  row["driver_id"].(string),
		})
	}

	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return vehicles, nil
}


