package config

import (
	"os"
)

//This is an interface to the file system. Several functions rely on IO with
// the OS file system. This interface allows use of a mock interface for
// unit testing functions without hitting our real file system.

// For a struct to implement this interface, it must implement all of the methods.
// Notice interface abstracts basic file operations - the things the functions
// I want to test are dependent upon/coupled to.
type FileSystem interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, permissions os.FileMode) error
	Getwd() (string, error)
}

// OSFilesystem is the real implementation. It uses the os package and represents the
// actual filesystem. It is an empty struct, taking up 0 bytes of mem. Signals to
// the compiler that we want to use its receivers without using memory. Comparable to
// but not the same as function overloading.
type OSFileSystem struct{}

// Receivers for OSFileSystem that wrap the standard os library IO functions.
// We call these in our functions as if they were the os implementations.
// Note these are all value receivers asd OSFileSystem is stateless - i.e 0 size.
func (OSFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// Similar wrapper for interface
func (OSFileSystem) WriteFile(filename string, data []byte, permissions os.FileMode) error {
	return os.WriteFile(filename, data, permissions)
}

func (OSFileSystem) Getwd() (string, error) {
	return os.Getwd()
}

// Mock filesystem - purely for injecting to tests so we can mock up files for
// unit test io. Note we include members here which will be accessed by receivers
// when mocking input.
type MockFileSystem struct {
	// files can be accessed by 'filepath' key, giving byte slice
	// as a normal file would when io.Read
	files map[string][]byte
	// A working directory we can set for testing
	wd string
}

// Pointer (as not 0 mem) receivers to MockFilesystem which will be called by the
// functions we want to test - though during normal use, an OSFileSystem and
//its corresponding receivers will be used.

func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	data, ok := m.files[filename]
	if !ok {
		// Need to flag that the file is not found.
		return nil, os.ErrNotExist
	}

	return data, nil
}

// Just add to the map representing our file system.
func (m *MockFileSystem) WriteFile(filename string, data []byte, permissions os.FileMode) error {
	m.files[filename] = data
	return nil
}

// Simply grab the wd from our MockFileSystem. We can set this directly when testing.
func (m *MockFileSystem) Getwd() (string, error) {
	return m.wd, nil
}
