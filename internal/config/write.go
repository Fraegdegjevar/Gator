package config

import (
	"encoding/json"
	"fmt"
)

// Write the config struct to JSON config file .gatorconfig.json
func write(fs FileSystem, conf *Config) error {
	// Get config filepath
	filePath, err := getConfigFilePath(fs)
	if err != nil {
		return fmt.Errorf("error writing config to %s: %w", configFileName, err)
	}

	//Marshal JSON
	data, err := json.MarshalIndent(conf, "", "	")
	if err != nil {
		return fmt.Errorf("error marshaling config to JSONL %w", err)
	}
	// Write file, Permission bits are linux permissions. First number, 0, tells Go
	// that this is an octal (base 8) number. Second number are owner permissions, third
	// group perms, fourth user perms. Permissions are read = 4, write = 2, exec = 1
	// So 6 = 4+2 = rw-, 4 = 4 = r--.
	// In other words -rw-r--r--
	//Note this is the WriteFile method for our filesystem - wraps os.WriteFile
	// if using an OSFileSystem - otherwise our test mocksystem WriteFile.
	err = fs.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write marshaled JSON to %v: %w", filePath, err)
	}
	return nil
}
