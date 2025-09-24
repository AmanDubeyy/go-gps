package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/joho/godotenv"
	"ondc/model"
)

var GpsChan = make(chan model.GPSData, 10000)
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

	for i := 0; i < workerCount; i++ {
		go InsertionWorker(i, repo)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/insert", InsertHandler())
	mux.HandleFunc("/route", SearchByID(repo))
	mux.HandleFunc("/search", SearchHandler(repo))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("🚍 GPS API running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	close(GpsChan)

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
}
