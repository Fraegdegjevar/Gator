package config

import (
	"encoding/json"
	"fmt"
	"io"
)

type Config struct {
	DB_URL string `json:"db_url"`
}

func Read(f io.Reader) (Config, error) {
	//Struct for unmarshaled JSON
	config := Config{}

	// Note newdecoder can take the file as an io.Reader
	// and both stream and decode into a struct
	err := json.NewDecoder(f).Decode(&config)

	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling JSON file: %w", err)
	}
	return config, nil
}
