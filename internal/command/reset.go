package command

import (
	"context"
	"fmt"

	"github.com/Fraegdegjevar/Gator/internal/config"
)

func HandlerReset(fs config.FileSystem, s *State, cmd Command) error {
	// No args required - just reset and confirm.
	err := s.Db.DeleteUsers(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("Successfully reset all users.")
	return nil
}
