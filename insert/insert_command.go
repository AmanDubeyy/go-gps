package insert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"ondc/model"
)

func generateGPSData(id int) model.GPSData {
	return model.GPSData{
		VehicleID:   strconv.Itoa(id),
		VehicleType: []string{"Car", "Bus", "Bike", "Truck"}[rand.Intn(4)],
		RouteID:     fmt.Sprintf("R-%d", rand.Intn(100)),
		City:        []string{"Mumbai", "Delhi", "Bangalore", "Chennai"}[rand.Intn(4)],
		Lat:         19.0 + rand.Float64(),  // ~19.x
		Lon:         72.0 + rand.Float64(),  // ~72.x
		Speed:       20 + rand.Float64()*80, // between 20–100
		RegNo:       []string{"HR10B7373", "HR10B7371", "HR10B7370", "HR10B7372"}[rand.Intn(4)],
	}
}

func InsertRecords(n int) error {
	client := &http.Client{Timeout: 5 * time.Second}
	url := "http://localhost:8080/insert"

	for i := 1; i <= n; i++ {
		data := generateGPSData(i)

		body, err := json.Marshal(data)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		fmt.Printf("Inserted record %d: %s\n", i, data.VehicleID)
		time.Sleep(10 * time.Millisecond) // small delay to avoid overwhelming
	}
	return nil
}
