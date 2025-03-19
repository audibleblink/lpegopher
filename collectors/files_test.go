package collectors

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitOutputFiles(t *testing.T) {
	t.Run(
		"Output files are created with appropriate writers initialized",
		func(t *testing.T) {
			// Create a temporary directory for the test
			testDir, err := os.MkdirTemp("", "lpegopher_test")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(testDir)

			// Change working directory to test directory
			origDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(origDir)
			os.Chdir(testDir)

			// Initialize output files
			InitOutputFiles()

			// Check if all files were created
			filesToCheck := []string{
				ExeFile,
				DllFile,
				DirFile,
				PrincipalFile,
				RelsFile,
				DepsFile,
				RunnersFile,
				ImportFile,
			}

			for _, file := range filesToCheck {
				path := filepath.Join(testDir, file)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", file)
				}
			}

			// Check if writers were initialized
			for file, writer := range writers {
				if writer == nil {
					t.Errorf("Writer for %s was not initialized", file)
				}
			}

			// Clean up
			FlushAndClose()
		},
	)
}

func TestFlashAndClose(t *testing.T) {
	t.Run("Data is properly flushed to disk when FlushAndClose is called", func(t *testing.T) {
		// Create a temporary directory for the test
		testDir, err := os.MkdirTemp("", "lpegopher_test")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(testDir)

		// Change working directory to test directory
		origDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current directory: %v", err)
		}
		defer os.Chdir(origDir)
		os.Chdir(testDir)

		// Initialize output files
		InitOutputFiles()

		// Write some test data to verify flush works
		if writers[ExeFile] != nil {
			writers[ExeFile].Write([]byte("test data"))
		}

		// Flush and close files
		FlushAndClose()

		// Check if data was flushed
		content, err := os.ReadFile(ExeFile)
		if err != nil {
			t.Errorf("Failed to read test file: %v", err)
		} else if string(content) != "test data" {
			t.Errorf("Expected 'test data', got '%s'", string(content))
		}
	})
}
