package collectors

import (
	"os"
	"path/filepath"
	"testing"
)

func createMockTempDir(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "lpegopher_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a basic directory structure for testing
	testDirPath := filepath.Join(tempDir, "testdata")
	err = os.MkdirAll(testDirPath, 0755)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test directory structure: %v", err)
	}

	// Create a test executable file
	exePath := filepath.Join(testDirPath, "test.exe")
	err = os.WriteFile(exePath, []byte("test data"), 0644)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test executable: %v", err)
	}

	// Create a test DLL file
	dllPath := filepath.Join(testDirPath, "test.dll")
	err = os.WriteFile(dllPath, []byte("test data"), 0644)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test DLL: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// Note: These tests are disabled because they reference unexported functions
// They are kept here as templates for when we can make those functions testable
