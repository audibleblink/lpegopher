package collectors

import (
	"bytes"
	"strings"
	"testing"
)

func TestINodeMethods(t *testing.T) {
	// Set up test INode
	inode := INode{
		Name:   "test.exe",
		Path:   "/path/to/test.exe",
		Parent: "/path/to",
		Type:   "exe",
		DACL: DACL{
			Owner: &Principal{Name: "TestOwner"},
			Group: &Principal{Name: "TestGroup"},
		},
	}

	t.Run("ID generates consistent hash", func(t *testing.T) {
		id1 := inode.ID()
		id2 := inode.ID()

		if id1 == "" {
			t.Error("ID should not be empty")
		}

		if id1 != id2 {
			t.Errorf("ID should be consistent: got %s and %s", id1, id2)
		}
	})

	t.Run("CacheKey returns the correct value", func(t *testing.T) {
		cacheKey := inode.CacheKey()
		if cacheKey != "/path/to/test.exe" {
			t.Errorf("Expected CacheKey to be '/path/to/test.exe', got '%s'", cacheKey)
		}
	})

	t.Run("ToCSV formats correctly", func(t *testing.T) {
		csv := inode.ToCSV()

		if !strings.Contains(csv, inode.ID()) {
			t.Error("CSV should contain the node ID")
		}

		if !strings.Contains(csv, "test.exe") {
			t.Error("CSV should contain the node name")
		}

		// Check CSV has expected format (ID,Name,Path,Parent,Owner,Group)
		fields := strings.Split(strings.TrimSpace(csv), ",")
		if len(fields) != 6 {
			t.Errorf("Expected 6 CSV fields, got %d", len(fields))
		}
	})

	t.Run("Write outputs data and returns ID", func(t *testing.T) {
		var buf bytes.Buffer
		id := inode.Write(&buf)

		if id == "" {
			t.Error("Write should return a non-empty ID")
		}

		if id != inode.ID() {
			t.Errorf("Write should return the node ID: expected %s, got %s", inode.ID(), id)
		}

		if buf.Len() == 0 {
			t.Error("Write should output data to the buffer")
		}
	})
}

func TestPrincipalMethods(t *testing.T) {
	// Set up test Principal
	principal := Principal{
		Name:  "TestUser",
		Group: "TestGroup",
		Type:  "user",
	}

	t.Run("ID generates consistent hash", func(t *testing.T) {
		id1 := principal.ID()
		id2 := principal.ID()

		if id1 == "" {
			t.Error("ID should not be empty")
		}

		if id1 != id2 {
			t.Errorf("ID should be consistent: got %s and %s", id1, id2)
		}
	})

	t.Run("CacheKey returns the correct value", func(t *testing.T) {
		cacheKey := principal.CacheKey()
		if cacheKey != "TestUser" {
			t.Errorf("Expected CacheKey to be 'TestUser', got '%s'", cacheKey)
		}
	})

	t.Run("ToCSV formats correctly", func(t *testing.T) {
		csv := principal.ToCSV()

		if !strings.Contains(csv, principal.ID()) {
			t.Error("CSV should contain the principal ID")
		}

		if !strings.Contains(csv, "testuser") {
			t.Error("CSV should contain the principal name (lowercase)")
		}

		// Check CSV has fields for ID, Name, Group, and Type
		fields := strings.Split(strings.TrimSpace(csv), ",")
		// The implementation uses a buffer of size 6 but only sets 4 fields
		// which means we get a CSV with 6 fields with last two being empty
		if len(fields) < 4 {
			t.Errorf("Expected at least 4 CSV fields, got %d", len(fields))
		}

		// Check required fields are present
		if fields[0] != principal.ID() {
			t.Errorf("First field should be ID, got %s", fields[0])
		}
		if !strings.Contains(fields[1], "testuser") {
			t.Errorf("Second field should contain name, got %s", fields[1])
		}
		if !strings.Contains(fields[2], "testgroup") {
			t.Errorf("Third field should contain group, got %s", fields[2])
		}
		if fields[3] != "user" {
			t.Errorf("Fourth field should be type, got %s", fields[3])
		}
	})

	t.Run("Write outputs data and returns ID", func(t *testing.T) {
		var buf bytes.Buffer
		id := principal.Write(&buf)

		if id == "" {
			t.Error("Write should return a non-empty ID")
		}

		if id != principal.ID() {
			t.Errorf(
				"Write should return the principal ID: expected %s, got %s",
				principal.ID(),
				id,
			)
		}

		if buf.Len() == 0 {
			t.Error("Write should output data to the buffer")
		}
	})
}

func TestRelMethods(t *testing.T) {
	// Set up test relationship
	rel := Rel{
		Start: "start123",
		Rel:   "CONTAINS",
		End:   "end456",
	}

	t.Run("ID generates consistent hash", func(t *testing.T) {
		id1 := rel.ID()
		id2 := rel.ID()

		if id1 == "" {
			t.Error("ID should not be empty")
		}

		if id1 != id2 {
			t.Errorf("ID should be consistent: got %s and %s", id1, id2)
		}
	})

	t.Run("CacheKey returns the correct value", func(t *testing.T) {
		cacheKey := rel.CacheKey()
		expectedCSV := "start123,CONTAINS,end456\n"
		if cacheKey != expectedCSV {
			t.Errorf("Expected CacheKey to be '%s', got '%s'", expectedCSV, cacheKey)
		}
	})

	t.Run("ToCSV formats correctly", func(t *testing.T) {
		csv := rel.ToCSV()

		if !strings.Contains(csv, "start123") {
			t.Error("CSV should contain the start ID")
		}

		if !strings.Contains(csv, "CONTAINS") {
			t.Error("CSV should contain the relationship type")
		}

		if !strings.Contains(csv, "end456") {
			t.Error("CSV should contain the end ID")
		}

		// Check CSV has expected format (Start,Rel,End)
		fields := strings.Split(strings.TrimSpace(csv), ",")
		if len(fields) != 3 {
			t.Errorf("Expected 3 CSV fields, got %d", len(fields))
		}
	})

	t.Run("Write outputs data and returns ID", func(t *testing.T) {
		var buf bytes.Buffer
		id := rel.Write(&buf)

		if id == "" {
			t.Error("Write should return a non-empty ID")
		}

		if id != rel.ID() {
			t.Errorf("Write should return the rel ID: expected %s, got %s", rel.ID(), id)
		}

		if buf.Len() == 0 {
			t.Error("Write should output data to the buffer")
		}
	})
}

func TestDepMethods(t *testing.T) {
	// Set up test dependency
	dep := Dep{
		Name: "kernel32.dll",
	}

	t.Run("ID generates consistent hash", func(t *testing.T) {
		id1 := dep.ID()
		id2 := dep.ID()

		if id1 == "" {
			t.Error("ID should not be empty")
		}

		if id1 != id2 {
			t.Errorf("ID should be consistent: got %s and %s", id1, id2)
		}
	})

	t.Run("CacheKey returns the correct value", func(t *testing.T) {
		cacheKey := dep.CacheKey()
		if cacheKey != "kernel32.dll" {
			t.Errorf("Expected CacheKey to be 'kernel32.dll', got '%s'", cacheKey)
		}
	})

	t.Run("ToCSV formats correctly", func(t *testing.T) {
		csv := dep.ToCSV()

		if !strings.Contains(csv, dep.ID()) {
			t.Error("CSV should contain the dependency ID")
		}

		if !strings.Contains(csv, "kernel32.dll") {
			t.Error("CSV should contain the dependency name")
		}

		// Check CSV has expected format (ID,Name)
		fields := strings.Split(strings.TrimSpace(csv), ",")
		if len(fields) != 2 {
			t.Errorf("Expected 2 CSV fields, got %d", len(fields))
		}
	})

	t.Run("Write outputs data and returns ID", func(t *testing.T) {
		var buf bytes.Buffer
		id := dep.Write(&buf)

		if id == "" {
			t.Error("Write should return a non-empty ID")
		}

		if id != dep.ID() {
			t.Errorf("Write should return the dep ID: expected %s, got %s", dep.ID(), id)
		}

		if buf.Len() == 0 {
			t.Error("Write should output data to the buffer")
		}
	})
}

func TestPERunnerMethods(t *testing.T) {
	// Set up test PERunner
	peRunner := PERunner{
		Name:     "TestService",
		Type:     "service",
		Args:     "-k test",
		RunLevel: "system",
		Exe: &INode{
			Name:   "service.exe",
			Path:   "/path/to/service.exe",
			Parent: "/path/to",
		},
		Context: &Principal{
			Name: "SYSTEM",
		},
	}

	t.Run("ID generates consistent hash", func(t *testing.T) {
		id1 := peRunner.ID()
		id2 := peRunner.ID()

		if id1 == "" {
			t.Error("ID should not be empty")
		}

		if id1 != id2 {
			t.Errorf("ID should be consistent: got %s and %s", id1, id2)
		}
	})

	t.Run("CacheKey returns the correct value", func(t *testing.T) {
		cacheKey := peRunner.CacheKey()
		if cacheKey != "TestService" {
			t.Errorf("Expected CacheKey to be 'TestService', got '%s'", cacheKey)
		}
	})

	t.Run("ToCSV formats correctly", func(t *testing.T) {
		csv := peRunner.ToCSV()

		if !strings.Contains(csv, peRunner.ID()) {
			t.Error("CSV should contain the runner ID")
		}

		if !strings.Contains(csv, "testservice") {
			t.Error("CSV should contain the runner name (lowercase)")
		}

		if !strings.Contains(csv, "service") {
			t.Error("CSV should contain the runner type")
		}

		// Check CSV has expected format with 8 fields
		fields := strings.Split(strings.TrimSpace(csv), ",")
		if len(fields) != 8 {
			t.Errorf("Expected 8 CSV fields, got %d", len(fields))
		}
	})

	t.Run("Write outputs data and returns ID", func(t *testing.T) {
		var buf bytes.Buffer
		id := peRunner.Write(&buf)

		if id == "" {
			t.Error("Write should return a non-empty ID")
		}

		if id != peRunner.ID() {
			t.Errorf("Write should return the runner ID: expected %s, got %s", peRunner.ID(), id)
		}

		if buf.Len() == 0 {
			t.Error("Write should output data to the buffer")
		}
	})
}

func TestWriteItems(t *testing.T) {
	// Create multiple principals
	principals := []Principal{
		{Name: "SYSTEM", Group: "NT AUTHORITY", Type: "user"},
		{Name: "Administrator", Group: "BUILTIN", Type: "user"},
		{Name: "Users", Group: "BUILTIN", Type: "group"},
	}

	// Test WriteItems
	var buf bytes.Buffer
	ids := WriteItems(principals, &buf)

	// Check that IDs are returned
	if len(ids) != len(principals) {
		t.Errorf("Expected %d IDs, got %d", len(principals), len(ids))
	}

	for _, id := range ids {
		if id == "" {
			t.Error("Expected non-empty ID to be returned")
		}
	}

	// Check that some data was written
	if buf.Len() == 0 {
		t.Error("Expected data to be written to buffer")
	}

	// Test writing again - should be cached
	buf.Reset()
	ids2 := WriteItems(principals, &buf)

	// IDs should be the same
	for i, id := range ids {
		if id != ids2[i] {
			t.Errorf("Expected same ID to be returned on second write, got %s != %s", id, ids2[i])
		}
	}

	// Buffer should be empty as items were cached
	if buf.Len() > 0 {
		t.Error("Expected no data to be written for cached items")
	}
}

func TestHashFor(t *testing.T) {
	t.Run("Returns consistent hashes", func(t *testing.T) {
		testString := "test string"
		hash1 := hashFor(testString)
		hash2 := hashFor(testString)

		if hash1 == "" {
			t.Error("Hash should not be empty")
		}

		if hash1 != hash2 {
			t.Errorf("Hash should be consistent for the same input: got %s and %s", hash1, hash2)
		}
	})

	t.Run("Different inputs yield different hashes", func(t *testing.T) {
		hash1 := hashFor("input1")
		hash2 := hashFor("input2")

		if hash1 == hash2 {
			t.Error("Different inputs should produce different hashes")
		}
	})

	t.Run("Paths are normalized", func(t *testing.T) {
		hash1 := hashFor(`C:\Windows\System32`)
		hash2 := hashFor(`c:\windows\system32`)

		if hash1 != hash2 {
			t.Error("Path case differences should be normalized")
		}
	})
}
