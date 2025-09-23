package main

import "os"

type Config struct {
	InfluxURL    string
	InfluxToken  string
	InfluxDBName string
}

func LoadConfig() Config {
	return Config{
		InfluxURL:    getEnv("INFLUX_URL", "http://localhost:8181"),
		InfluxToken:  getEnv("INFLUX_TOKEN", "my-secret-token"),
		InfluxDBName: getEnv("INFLUX_DB", "gpszz"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
