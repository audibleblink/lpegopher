package processor

import (
	"testing"
)

func TestDataPrefix(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Empty URL",
			url:      "",
			expected: "file://",
		},
		{
			name:     "With URL",
			url:      "localhost:8080",
			expected: "http://localhost:8080",
		},
		{
			name:     "With IP address",
			url:      "192.168.1.1:8080",
			expected: "http://192.168.1.1:8080",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := dataPrefix(tc.url)
			if result != tc.expected {
				t.Errorf("dataPrefix(%q) = %q, expected %q", tc.url, result, tc.expected)
			}
		})
	}
}