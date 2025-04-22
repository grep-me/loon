package load

import (
	"os"
	"path/filepath"
)

// GetExecutables scans a given directory and returns a slice of paths of executable files.
func GetExecutables(dirPath string) ([]string, error) {
	// Create a slice to store executable file paths
	var executablePaths []string

	// Walk through the directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is not a directory and if it's executable
		if !info.IsDir() && isExecutable(path) {
			executablePaths = append(executablePaths, "./"+path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return executablePaths, nil
}

// isExecutable checks if the file is executable.
func isExecutable(path string) bool {
	// Get the file information
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check the file permissions and ensure it's executable
	mode := info.Mode()
	return mode&0111 != 0 // The 0111 check ensures that any of the executable bits are set
}
