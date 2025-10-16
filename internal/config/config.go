package config

import (
	"errors"
	"fmt"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DB_URL            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func getConfigFilePath(fs FileSystem) (string, error) {
	// Ensure file exists and read it in
	filePath, err := fs.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting %s filepath: %w", configFileName, err)
	}
	filePath += configFileName
	return filePath, nil
}

func (c *Config) SetUser(fs FileSystem, username string) error {
	if len(username) < 1 {
		return errors.New("no username supplied to SetUser")
	}

	c.Current_user_name = username

	err := write(fs, c)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	return nil
}
