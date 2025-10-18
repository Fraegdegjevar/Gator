package command

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Fraegdegjevar/Gator/internal/config"
)

func TestRun(t *testing.T) {

	// Track if command was called or not. This is rebound per test and checked locally per test.
	var commandCalled bool

	cmds := Commands{
		Registry: make(map[string]func(config.FileSystem, *State, Command) error),
	}

	s := &State{}

	fs := config.OSFileSystem{}
	// Sentinel/mock error for test
	errMissingMockArgs := errors.New("must suppy an argument to handlerMock")

	// Register a mock command - we just want to check it is run if it needs to be run,
	// or that it errors if we do not get the expected args passed in
	handlerMock := func(fs config.FileSystem, s *State, cmd Command) error {
		commandCalled = true
		if len(cmd.Args) < 1 {
			return errMissingMockArgs
		}
		return nil
	}
	cmds.Register("mock", handlerMock)

	cases := []struct {
		name           string
		inputCommand   Command
		expectedCalled bool
		expectedError  error
	}{
		{
			name:           "successful run",
			inputCommand:   Command{Name: "mock", Args: []string{"arg1"}},
			expectedCalled: true,
			expectedError:  nil,
		},
		{
			name:           "missing args for command",
			inputCommand:   Command{Name: "mock", Args: []string{}},
			expectedCalled: true,
			expectedError:  errMissingMockArgs,
		},
		{
			name:           "command not found",
			inputCommand:   Command{Name: "noexist", Args: []string{"somearg"}},
			expectedCalled: false,
			expectedError:  ErrCommandNotFound,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			commandCalled = false

			err := cmds.Run(fs, s, tt.inputCommand)
			// Is the right error (or nil) returned?
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error: %v got: %v", tt.expectedError, err)
			}

			// Did we call/notcall command as expected?
			if tt.expectedCalled != commandCalled {
				t.Errorf("expected command to be called: %v, got: %v", tt.expectedCalled, commandCalled)
			}
		})
	}
}

func TestHandlerLogin(t *testing.T) {

	// Define a pre-existing config file content for tests that need it.
	preExistingConfigContent := `{"db_url":"","current_user_name":"default"}`

	cases := []struct {
		name           string
		filesystem     *config.MockFileSystem
		state          *State
		cmd            Command
		expectedError  error
		expectedConfig config.Config
	}{
		{
			name: "success",
			filesystem: &config.MockFileSystem{
				Wd: "test",
				Files: map[string][]byte{
					"test/.gatorconfig.json": []byte(preExistingConfigContent),
				},
			},
			state:          &State{Config: &config.Config{DBURL: "", CurrentUserName: "default"}},
			cmd:            Command{Args: []string{"testuser"}},
			expectedError:  nil,
			expectedConfig: config.Config{DBURL: "", CurrentUserName: "testuser"},
		},
		{
			name: "fail no username",
			filesystem: &config.MockFileSystem{
				Wd: "test",
				Files: map[string][]byte{
					"test/.gatorconfig.json": []byte(preExistingConfigContent),
				},
			},
			state:          &State{Config: &config.Config{DBURL: "", CurrentUserName: "default"}},
			cmd:            Command{Args: []string{}},
			expectedError:  config.ErrNoUsername,
			expectedConfig: config.Config{DBURL: "", CurrentUserName: "default"},
		},
		{
			name: "fail as SetUser fails",
			filesystem: &config.MockFileSystem{
				Wd: "test",
				Files: map[string][]byte{
					"test/.gatorconfig.json": []byte(preExistingConfigContent),
				},
				WriteFileShouldError: config.ErrWriteFail,
			},
			state:          &State{&config.Config{DBURL: "", CurrentUserName: "default"}},
			cmd:            Command{Args: []string{"testfail"}},
			expectedError:  config.ErrWriteFail,
			expectedConfig: config.Config{DBURL: "", CurrentUserName: "testfail"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			err := HandlerLogin(tt.filesystem, tt.state, tt.cmd)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error: %v but got: %v", tt.expectedError, err)
			}

			if !reflect.DeepEqual(tt.expectedConfig, *tt.state.Config) {
				t.Errorf("Expected config: %v got: %v", tt.expectedConfig, *tt.state.Config)
			}

		})
	}
}
