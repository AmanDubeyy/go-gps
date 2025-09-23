package main

import (
	"context"
	"time"
	"fmt"

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
		"bus_locations",
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

	fmt.Print("here")

	return r.client.WritePoints(ctx, []*influxdb3.Point{point})
}

func (r *InfluxGPSRepository) GetRoute(vehicleID string, start, end time.Time) ([]GPSData, error) {
	ctx := context.Background()

	query := fmt.Sprintf(`
	from(bucket:"%s")
	  |> range(start: %s, stop: %s)
	  |> filter(fn: (r) => r._measurement == "bus_locations")
	  |> filter(fn: (r) => r.vehicle_id == "%s")
	  |> sort(columns:["_time"])
	`, r.bucket, start.Format(time.RFC3339), end.Format(time.RFC3339), vehicleID)

	result, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var vehicles []GPSData
	for result.Next() {
		row := result.Value()
		vehicles = append(vehicles, GPSData{
        VehicleID: row["vehicle_id"].(string),
        RouteID:   row["route_id"].(string),
        City:      row["city"].(string),
        Status:    row["status"].(string),
        Lat:       row["lat"].(float64),
        Lon:       row["lon"].(float64),
        Speed:     row["speed"].(float64),
        Heading:   int(row["heading"].(float64)),
        DriverID:  row["driver_id"].(string),
    	})
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return vehicles, nil
}

func (r *InfluxGPSRepository) SearchVehicles(minLat, maxLat, minLon, maxLon float64) ([]GPSData, error) {
	ctx := context.Background()

	// query := fmt.Sprintf(`
	// 	SELECT *
	// 	FROM "bus_locations"
	// 	WHERE lat >= %f AND lat <= %f
	// 	AND lon >= %f AND lon <= %f
	// 	AND time >= NOW() - INTERVAL 1 HOUR
	// 	`, minLat, maxLat, minLon, maxLon,
	// )

	query := "SELECT * FROM bus_locations"

	iter, err := r.client.Query(ctx, query)
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


