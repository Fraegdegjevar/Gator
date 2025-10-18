package config

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {

	// Our mock filesystem with the file
	fsCorrect := MockFileSystem{
		Wd:    "test",
		Files: make(map[string][]byte),
	}
	// Mock json we insert into mockfs to test read
	mockConfig := Config{
		DBURL:           "testURL",
		CurrentUserName: "testuser",
	}

	marshaledJSON, err := json.MarshalIndent(mockConfig, "", "	")
	if err != nil {
		t.Fatalf("failed to marshal JSON when setting up for tests")
	}
	fsCorrect.Files["test/.gatorconfig.json"] = marshaledJSON

	//Mock filesystem that will be missing the config
	// Our mock filesystem with the file
	fsMissing := MockFileSystem{
		Wd:    "test",
		Files: make(map[string][]byte),
	}

	cases := []struct {
		name            string
		inputFileSystem FileSystem
		expectedConfig  Config
		expectedErr     bool
	}{
		{
			name:            "correct json",
			inputFileSystem: &fsCorrect,
			expectedConfig:  mockConfig,
			expectedErr:     false,
		},
		{
			name:            "config file missing",
			inputFileSystem: &fsMissing,
			expectedConfig:  Config{},
			expectedErr:     true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			conf, err := Read(tt.inputFileSystem)

			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v but got: %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(conf, tt.expectedConfig) {
				t.Errorf("expected config: %v, got: %v", tt.expectedConfig, conf)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	fs := &MockFileSystem{
		Files: make(map[string][]byte),
		Wd:    "test",
	}

	correctConfig := Config{
		DBURL:           "testurl",
		CurrentUserName: "testuser",
	}

	blankConfig := Config{}

	cases := []struct {
		name           string
		inputConfig    *Config
		expectedConfig *Config
		expectedError  bool
	}{
		{
			name:           "correct populated config fields",
			inputConfig:    &correctConfig,
			expectedConfig: &correctConfig,
		},
		{
			name:           "blank config",
			inputConfig:    &blankConfig,
			expectedConfig: &blankConfig,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := write(fs, tt.inputConfig)
			if err != nil {
				t.Fatalf("error writing config for testing: %v", err)
			}

			actualConf, err := Read(fs)
			if err != nil {
				t.Fatalf("error reading the written config file: %v", err)
			}

			if !reflect.DeepEqual(*tt.expectedConfig, actualConf) {
				t.Errorf("expected config: %v, got: %v", *tt.expectedConfig, actualConf)
			}

		})
	}
}

func TestSetUser(t *testing.T) {
	//Test where success -user supplied, our config should be updated & write called
	//Test where no username supplied - don't run, just error ErrNoUsername
	//Test where username suppplied but write error - run, config updated, but handles error if file not written

	// PAss in a config with a 'default' username - to ensure our mockFilesystem has a username set.
	//conf := &Config{CurrentUserName: "default"}

	// each shadowed tt test iteration. All functions inside t.Run iterations will need to be
	// pased a pointer to their copy of the filesystem.
	// write is already tested - marshals and writes etc.

	cases := []struct {
		name                   string
		inputUsername          string
		inputConfig            *Config
		fileSystem             MockFileSystem
		expectedError          error
		expectedConfigUserName string
		expectedWrite          bool
		expectedFileUserName   string
	}{
		{
			name:                   "no username error and no write_tst",
			inputUsername:          "",
			inputConfig:            &Config{CurrentUserName: "default"},
			fileSystem:             MockFileSystem{Wd: "tst", Files: make(map[string][]byte), WriteCalled: 0},
			expectedError:          ErrNoUsername,
			expectedConfigUserName: "default",
			expectedWrite:          false,
			expectedFileUserName:   "default",
		},
		{
			name:                   "write_failed_tst",
			inputUsername:          "writefail",
			inputConfig:            &Config{},
			fileSystem:             MockFileSystem{Wd: "tst", WriteCalled: 0, Files: make(map[string][]byte)},
			expectedError:          ErrWriteFail,
			expectedConfigUserName: "writefail",
			expectedWrite:          true,
			expectedFileUserName:   "default",
		},
		{
			name:                   "success_set user_&_write_file_tst",
			inputUsername:          "testsuccess",
			inputConfig:            &Config{CurrentUserName: "default"},
			fileSystem:             MockFileSystem{Wd: "tst", Files: make(map[string][]byte), WriteCalled: 0},
			expectedError:          nil,
			expectedConfigUserName: "testsuccess",
			expectedWrite:          true,
			expectedFileUserName:   "testsuccess",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			// shadowing tt so we get a fresh tt for THIS TEST every iteration.
			tt := tt

			//Write a config to filesystem before trying to set user so we have something to go off!
			write(&tt.fileSystem, &Config{CurrentUserName: "default"})

			// If our expectedError is a write fail, set Filesystem to fail writes:
			if errors.Is(tt.expectedError, ErrWriteFail) {
				tt.fileSystem.WriteFileShouldError = ErrWriteFail
			}

			// check writeCount now to see if it is called below
			prevWriteCount := tt.fileSystem.WriteCalled

			err := tt.inputConfig.SetUser(&tt.fileSystem, tt.inputUsername)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			if tt.inputConfig.CurrentUserName != tt.expectedConfigUserName {
				t.Errorf("expected config username to be: %v, got: %v", tt.expectedConfigUserName, tt.inputConfig.CurrentUserName)
			}
			// Check write was called.
			actualWrite := (tt.fileSystem.WriteCalled != prevWriteCount)
			if tt.expectedWrite != actualWrite {
				t.Errorf("expected write to be called: %v, but got: %v", tt.expectedWrite, actualWrite)
			}

			//check the file's written username is correct by reading it using previously tested functions.
			// UNLESS we expect writes to fail - in which case
			readConf, err := Read(&tt.fileSystem)
			if err != nil {
				t.Fatalf("failed to read config from stored files in filesystem during test: %v, %v:", tt.name, err)
			}

			if readConf.CurrentUserName != tt.expectedFileUserName {
				t.Errorf("expected username stored in file to be: %v, got: %v", tt.expectedFileUserName, readConf.CurrentUserName)
			}

		})
	}
}
