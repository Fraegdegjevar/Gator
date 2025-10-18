package command

import (
	"errors"

	"github.com/Fraegdegjevar/Gator/internal/config"
)

//export a command not found error for use elsewhere + for testing
var ErrCommandNotFound = errors.New("command not found")

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Registry map[string]func(config.FileSystem, *State, Command) error
}

func (c *Commands) Run(fs config.FileSystem, s *State, cmd Command) error {
	fn, ok := c.Registry[cmd.Name]
	if !ok {
		return ErrCommandNotFound
	}

	err := fn(fs, s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(config.FileSystem, *State, Command) error) {
	c.Registry[name] = f
}
