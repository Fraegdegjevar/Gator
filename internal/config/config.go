package config

import (
	"errors"
	"fmt"
)

const configFileName = ".gatorconfig.json"

// exported error
var ErrNoUsername = errors.New("no username supplied")

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath(fs FileSystem) (string, error) {
	// Ensure file exists and read it in
	filePath, err := fs.GetUserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting %s filepath: %w", configFileName, err)
	}
	filePath += "/" + configFileName
	return filePath, nil
}

func (c *Config) SetUser(fs FileSystem, username string) error {
	if len(username) < 1 {
		return ErrNoUsername
	}

	c.CurrentUserName = username

	err := write(fs, c)
	if err != nil {
		return fmt.Errorf("error setting user in configuration file: %w", err)
	}

	fmt.Printf("user has been set: %v\n", username)
	return nil
}
