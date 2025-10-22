package util

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"HELLO", "hello"},
		{"Hello World", "hello world"},
		{"", ""},
	}

	for _, test := range tests {
		result := Lower(test.input)
		if result != test.expected {
			t.Errorf("Lower(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestPathFix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`C:\path\to\file.txt`, "c:/path/to/file.txt"},
		{`"C:\path\to\file.txt"`, "c:/path/to/file.txt"},
		{`C:\path,to\file.txt`, "c:/path.to/file.txt"},
		{` C:\path\to\file.txt `, "c:/path/to/file.txt"},
	}

	for _, test := range tests {
		result := PathFix(test.input)
		if result != test.expected {
			t.Errorf("PathFix(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestResolveEnvPath(t *testing.T) {
	// Save original environment and restore after test
	origEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range origEnv {
			parts := strings.SplitN(e, "=", 2)
			os.Setenv(parts[0], parts[1])
		}
	}()

	// Setup test environment
	os.Clearenv()
	os.Setenv("TEMP", "C:\\Temp")
	os.Setenv("PROGRAMFILES", "C:\\Program Files")

	tests := []struct {
		input    string
		expected string
	}{
		{"%TEMP%\\file.txt", "C:\\Temp\\file.txt"},
		{"%PROGRAMFILES%\\App\\data.txt", "C:\\Program Files\\App\\data.txt"},
		{"%NONEXISTENT%\\file.txt", "%NONEXISTENT%\\file.txt"},
		{"C:\\normal\\path.txt", "C:\\normal\\path.txt"},
		{"%TEMP%path.txt", "%TEMP%path.txt"},
		{"%TEMP%/file.txt", "%TEMP%/file.txt"},
	}

	for _, test := range tests {
		result := resolveEnvPath(test.input)
		if result != test.expected {
			t.Errorf("resolveEnvPath(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestEvaluatePath(t *testing.T) {
	// Save original environment and restore after test
	origEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, e := range origEnv {
			parts := strings.SplitN(e, "=", 2)
			os.Setenv(parts[0], parts[1])
		}
	}()

	// Setup test environment
	os.Clearenv()
	os.Setenv("TEMP", "C:\\Temp")
	os.Setenv("PROGRAMFILES", "C:\\Program Files")

	tests := []struct {
		input    string
		expected string
	}{
		{"%TEMP%\\file.txt", "C:\\Temp\\file.txt"},
		{"%PROGRAMFILES%\\App\\data.txt", "C:\\Program Files\\App\\data.txt"},
		{"%NONEXISTENT%\\file.txt", "%NONEXISTENT%\\file.txt"},
		{"C:\\normal\\path.txt", "C:\\normal\\path.txt"},
		{"%TEMP%path.txt", "%TEMP%path.txt"},
		{"%TEMP%/file.txt", "%TEMP%/file.txt"},
	}

	for _, test := range tests {
		result := EvaluatePath(test.input)
		if result != test.expected {
			t.Errorf("EvaluatePath(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestLineCount(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"Single line", 0},
		{"Line 1\nLine 2", 1},
		{"Line 1\nLine 2\nLine 3\n", 3},
		{"Line 1\r\nLine 2\r\nLine 3\r\n", 3},
	}

	for _, test := range tests {
		reader := bytes.NewBufferString(test.input)
		count, err := LineCount(reader)
		if err != nil {
			t.Errorf("LineCount(%q) returned error: %v", test.input, err)
		}
		if count != test.expected {
			t.Errorf("LineCount(%q) = %d, expected %d", test.input, count, test.expected)
		}
	}
}

func TestRand(t *testing.T) {
	// Tests with the existing source
	result1 := Rand()
	if len(result1) != 6 {
		t.Errorf("Rand() returned string of length %d, expected 6", len(result1))
	}

	// Check that only valid characters are used
	for _, c := range result1 {
		if !strings.ContainsRune(letterBytes, c) {
			t.Errorf("Rand() returned string containing invalid character %q", c)
		}
	}
}

func TestSmoothBrainPath(t *testing.T) {
	tests := []struct {
		cmdline      string
		expectedBin  string
		expectedArgs string
	}{
		{
			`"C:\Program Files\App\app.exe" --config=test.json`,
			"C:\\Program Files\\App\\app.exe",
			"--config=test.json",
		},
		{
			`C:\Windows\System32\cmd.exe /c echo hello`,
			"C:\\Windows\\System32\\cmd.exe",
			"/c echo hello",
		},
		{
			`C:\Windows\rundll32.exe shell32.dll,Control_RunDLL`,
			"C:\\Windows\\rundll32.exe",
			"shell32.dll,Control_RunDLL",
		},
		{`program.exe arg1 arg2`, "program.exe", "arg1 arg2"},
	}

	for _, test := range tests {
		bin, args := SmoothBrainPath(test.cmdline)
		if bin != test.expectedBin || args != test.expectedArgs {
			t.Errorf("SmoothBrainPath(%q) = (%q, %q), expected (%q, %q)",
				test.cmdline, bin, args, test.expectedBin, test.expectedArgs)
		}
	}
}
