package main

import (
	"context"
	"fmt"
)

func InsertionWorker(id int, repo *InfluxGPSRepository) {
	for data := range GpsChan {
		if err := repo.InsertGPS(context.Background(), data); err != nil {
			fmt.Printf("[Worker %d] failed to insert GPS: %v", id, err)
		}
	}
}