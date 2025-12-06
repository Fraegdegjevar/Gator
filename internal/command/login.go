package command

import (
	"context"
	"fmt"
	"os"

	"github.com/Fraegdegjevar/Gator/internal/config"
)

// Login updates the config current user.
// but also checks that the user exists in the database
// before logging in.
func HandlerLogin(fs config.FileSystem, s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return config.ErrNoUsername
	}
	username := cmd.Args[0]

	//Is user in db? Exit with coe 1 if not.
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("error getting supplied user from database: %v\n", err)
		os.Exit(1)
	}

	err = s.Config.SetUser(fs, username)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	return nil
}
