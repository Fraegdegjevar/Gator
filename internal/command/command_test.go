package command

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/Fraegdegjevar/Gator/internal/config"
	"github.com/Fraegdegjevar/Gator/internal/database"
	_ "github.com/lib/pq"
)

// Define a pre-existing config file stub with test database details for tests that need it.
var testConfigContent = `{"db_url":"postgres://postgres:postgres@localhost:5432/gator_test?sslmode=disable","current_user_name":""}`

// A config struct we can test against - using the magic values stored in preExistingConfigContent
var testConfig config.Config

// Go runtime executes init function before main or test functions.
// init passes the db_url magic value (set this above based on your test database) into a config struct to allow for easy reference
// in tests below.
func init() {
	err := json.Unmarshal([]byte(testConfigContent), &testConfig)
	if err != nil {
		panic("failed to unmarshal testConfigContent byte(string): " + err.Error())
	}
}

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

	//Note that the Db in state is a database *Queries object using .New() on a SQL database connection.
	// Test by hitting a test database
	cases := []struct {
		name           string
		filesystem     *config.FakeFileSystem
		state          *State
		cmd            Command
		expectedError  error
		expectedConfig config.Config
	}{
		{
			name: "success",
			filesystem: &config.FakeFileSystem{
				Homedir: "test",
				Files: map[string][]byte{
					"test/.gatorconfig.json": []byte(testConfigContent),
				},
			},
			state:          &State{Config: &config.Config{}},
			cmd:            Command{Args: []string{"testuser"}},
			expectedError:  nil,
			expectedConfig: config.Config{DBURL: testConfig.DBURL, CurrentUserName: "testuser"},
		},
		{
			name: "fail no username",
			filesystem: &config.FakeFileSystem{
				Homedir: "test",
				Files: map[string][]byte{
					"test/.gatorconfig.json": []byte(testConfigContent),
				},
			},
			state:          &State{Config: &config.Config{}},
			cmd:            Command{Args: []string{}},
			expectedError:  config.ErrNoUsername,
			expectedConfig: config.Config{DBURL: testConfig.DBURL, CurrentUserName: "default"},
		},
		{
			name: "fail as SetUser fails",
			filesystem: &config.FakeFileSystem{
				Homedir: "test",
				Files: map[string][]byte{
					"test/.gatorconfig.json": []byte(testConfigContent),
				},
				WriteFileShouldError: config.ErrWriteFail,
			},
			state:          &State{Config: &config.Config{}},
			cmd:            Command{Args: []string{"testfail"}},
			expectedError:  config.ErrWriteFail,
			expectedConfig: config.Config{DBURL: testConfig.DBURL, CurrentUserName: "testfail"},
		},
		{
			name: "fail as user does not exist in db",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			// Open a connection pool to the test database specified in testConfig initialised with init()
			tdb, err := sql.Open("postgres", testConfig.DBURL)
			if err != nil {
				t.Fatalf("error connecting to test database: %v", err)
			}
			tt.state.Db = database.New(tdb)
			err = HandlerLogin(tt.filesystem, tt.state, tt.cmd)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error: %v but got: %v", tt.expectedError, err)
			}

			if !reflect.DeepEqual(tt.expectedConfig, *tt.state.Config) {
				t.Errorf("Expected config: %v got: %v", tt.expectedConfig, *tt.state.Config)
			}

		})
	}
}

func TestHandlerRegister(t *testing.T) {

}
