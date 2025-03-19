package collectors

import (
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"github.com/minio/highwayhash"

	"github.com/audibleblink/lpegopher/util"
)

// TestHasherCompatibility ensures the new implementation produces the same
// hash values as the old implementation
func TestHasherCompatibility(t *testing.T) {
	testCases := []string{
		"simple string",
		"C:\\Windows\\System32\\notepad.exe",
		"c:\\windows\\system32\\NOTEPAD.EXE", // Should normalize to same as above
		"/usr/local/bin/bash",
		"This is a much longer string that would test performance with realistic data sizes that might be seen in normal usage scenarios",
	}

	for _, tc := range testCases {
		// Get hash using new implementation
		newHash, err := HashWithOptions(tc, true)
		if err != nil {
			t.Errorf("New hasher failed for input %q: %v", tc, err)
			continue
		}

		// Simulate old implementation for comparison
		oldHash := simulateOldHashImplementation(tc)

		// Verify both implementations produce the same hash
		if newHash != oldHash {
			t.Errorf("Hash mismatch for input %q:\nOld: %s\nNew: %s", tc, oldHash, newHash)
		}
	}
}

// simulateOldHashImplementation recreates the old hashing logic for comparison
func simulateOldHashImplementation(data string) string {
	data = util.PathFix(data)
	hash, _ := highwayhash.New(key)

	txt := strings.NewReader(data)
	io.Copy(hash, txt)

	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum)
}

// TestHasherOptions verifies that the optional parameters work correctly
func TestHasherOptions(t *testing.T) {
	// Test with and without normalization
	input := "C:\\Windows\\System32"
	inputLower := "c:\\windows\\system32"

	// With normalization, these should hash to the same value
	hashWithNorm, _ := HashWithOptions(input, true)
	hashLowerWithNorm, _ := HashWithOptions(inputLower, true)

	if hashWithNorm != hashLowerWithNorm {
		t.Errorf("Expected same hash with normalization: %s != %s", hashWithNorm, hashLowerWithNorm)
	}

	// Without normalization, these should hash to different values
	hashWithoutNorm, _ := HashWithOptions(input, false)
	hashLowerWithoutNorm, _ := HashWithOptions(inputLower, false)

	if hashWithoutNorm == hashLowerWithoutNorm {
		t.Errorf("Expected different hashes without normalization: both %s", hashWithoutNorm)
	}
}

// Benchmarks to compare old vs new implementation

func BenchmarkOldHashImplementation(b *testing.B) {
	data := "C:\\Windows\\System32\\notepad.exe"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		simulateOldHashImplementation(data)
	}
}

func BenchmarkNewHashImplementation(b *testing.B) {
	data := "C:\\Windows\\System32\\notepad.exe"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashFor(data)
	}
}

func BenchmarkNewHashWithOptions(b *testing.B) {
	data := "C:\\Windows\\System32\\notepad.exe"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		HashWithOptions(data, true)
	}
}

// Benchmark with different data sizes

func BenchmarkNewHashSmallString(b *testing.B) {
	data := "small"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashFor(data)
	}
}

func BenchmarkNewHashMediumString(b *testing.B) {
	data := "C:\\Program Files\\Common Files\\Microsoft Shared\\ClickToRun\\Updates\\16.0.16130.20218"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashFor(data)
	}
}

func BenchmarkNewHashLargeString(b *testing.B) {
	// Create a large string (approximately 1KB)
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		sb.WriteString(
			"This is a long string that we're using to test the performance of the hasher with larger inputs. ",
		)
	}
	data := sb.String()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashFor(data)
	}
}

// Benchmark parallel usage to test pool behavior

func BenchmarkNewHashParallel(b *testing.B) {
	data := "C:\\Windows\\System32\\notepad.exe"
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hashFor(data)
		}
	})
}
