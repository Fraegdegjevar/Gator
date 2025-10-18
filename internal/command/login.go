package command

import (
	"fmt"

	"github.com/Fraegdegjevar/Gator/internal/config"
)

func HandlerLogin(fs config.FileSystem, s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return config.ErrNoUsername
	}
	username := cmd.Args[0]
	err := s.Config.SetUser(fs, username)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	fmt.Printf("user has been set: %v\n", username)
	return nil
}
