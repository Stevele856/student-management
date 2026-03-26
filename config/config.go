package config

import "os"


type Config struct {
	Port string
	DataFile string
}

func Load() *Config{
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == ""{
		dataFile = "static/students.json"
	}

	return &Config{
		Port: port,
		DataFile: dataFile,
	}
}