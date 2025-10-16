package config

import (
	"encoding/json"
	"fmt"
)

// Note we don't call os.ReadFile etc. We call our filesystem interface.
func Read(fs FileSystem) (Config, error) {
	filePath, err := getConfigFilePath(fs)
	if err != nil {
		return Config{}, err
	}

	file, err := fs.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file %v: %w", filePath, err)
	}

	//Struct for unmarshaled JSON
	config := Config{}

	// Note newdecoder can take the file as an io.Reader
	// and both stream and decode into a struct
	err = json.Unmarshal(file, &config)

	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling JSON file: %w", err)
	}
	return config, nil
}
