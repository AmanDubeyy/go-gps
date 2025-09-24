package main

import "fmt"

func getByRouteID(vehicleID string) string {
	return fmt.Sprintf(`
		SELECT *
		FROM %s
		WHERE vehicle_id = '%s'
		ORDER BY time ASC
	`, VEHICLE_LOCATION, vehicleID)
}
