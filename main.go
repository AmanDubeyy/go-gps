package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/joho/godotenv"
)

var GpsChan = make(chan GPSData, 10000)
const workerCount = 100

func main() {
	godotenv.Load()
	cfg := LoadConfig()

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     cfg.InfluxURL,
		Token:    cfg.InfluxToken,
		Database: cfg.InfluxDBName,
	})
	if err != nil {
		log.Fatalf("failed to connect to influxdb: %v", err)
	}
	defer client.Close()

	repo := NewInfluxGPSRepository(client, cfg.InfluxDBName)

	for i := range workerCount{
		go InsertionWorker(i, repo)
	}

	http.HandleFunc("/insert", InsertHandler())
	http.HandleFunc("/route", RouteHandler(repo))
	http.HandleFunc("/search", SearchHandler(repo))


	fmt.Println("🚍 GPS API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

