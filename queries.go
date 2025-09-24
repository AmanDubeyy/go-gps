package main

import "fmt"

var (
	getAll = fmt.Sprintf("SELECT * FROM %s", VEHICLE_LOCATION)
)

func getByRouteID(vehicleID string) string {
	return fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE vehicle_id = '%s'
		ORDER BY time ASC
	`, VEHICLE_LOCATION, vehicleID)
}
