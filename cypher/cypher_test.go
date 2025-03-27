package cypher

import (
	"strings"
	"testing"

	"github.com/audibleblink/logerr"
)

// Create a mock version of neo4j.Driver that doesn't need a real database connection
type MockDriver struct{}

// Mock implementation for tests
func createTestQuery() *Query {
	b := &strings.Builder{}
	l := logerr.DefaultLogger()
	return &Query{
		b: b,
		d: nil, // We're not executing queries in tests
		l: l,
	}
}

func TestMerge(t *testing.T) {
	q := createTestQuery()

	q.Merge("n", "Node", "id", "123")

	expected := "MERGE (n:Node { id: '123' }) "
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestCreate(t *testing.T) {
	q := createTestQuery()

	q.Create("n", "Node", "id", "123")

	expected := "CREATE (n:Node { id: '123' }) "
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestMatch(t *testing.T) {
	q := createTestQuery()

	q.Match("n", "Node", "id", "123")

	expected := "MATCH (n:Node { id: '123' }) "
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestChaining(t *testing.T) {
	q := createTestQuery()

	q.Match("n", "Node", "id", "123").
		With("n").
		Match("m", "Related", "id", "456").
		Relate("n", "CONNECTS_TO", "m").
		Return()

	expected := "MATCH (n:Node { id: '123' }) WITH n\nMATCH (m:Related { id: '456' }) MERGE (n)-[:CONNECTS_TO]->(m) RETURN count(*)\n"
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestAppend(t *testing.T) {
	q := createTestQuery()

	q.Append("CUSTOM CYPHER")

	expected := "CUSTOM CYPHER\n"
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestWith(t *testing.T) {
	q := createTestQuery()

	q.With("n, count(r) as count")

	expected := "WITH n, count(r) as count\n"
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestEndMerge(t *testing.T) {
	q := createTestQuery()

	q.EndMerge()

	expected := "WITH count(*) as dummy\n"
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestReturn(t *testing.T) {
	q := createTestQuery()

	q.Return()

	expected := "RETURN count(*)\n"
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestTerminate(t *testing.T) {
	q := createTestQuery()

	q.Terminate()

	expected := "\n"
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestSet(t *testing.T) {
	q := createTestQuery()

	props := map[string]string{
		"name": "test",
		"age":  "30",
	}

	q.Set("n", props)

	// Since map iteration order is not guaranteed, check for both possible outputs
	expected1 := "SET n.name = 'test', n.age = '30' "
	expected2 := "SET n.age = '30', n.name = 'test' "

	result := q.String()
	if result != expected1 && result != expected2 {
		t.Errorf("Expected query %q or %q, got %q", expected1, expected2, result)
	}
}

func TestRelate(t *testing.T) {
	q := createTestQuery()

	q.Relate("n", "KNOWS", "m")

	expected := "MERGE (n)-[:KNOWS]->(m) "
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestRaw(t *testing.T) {
	q := createTestQuery()

	// Add some initial content and then replace it
	q.Append("INITIAL CONTENT")
	q.Raw("REPLACED CONTENT")

	expected := "REPLACED CONTENT"
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestReset(t *testing.T) {
	q := createTestQuery()

	// Add some content and then reset
	q.Append("SOME CONTENT")
	q.Reset()

	expected := ""
	if q.String() != expected {
		t.Errorf("Expected empty query, got %q", q.String())
	}
}

func TestString(t *testing.T) {
	q := createTestQuery()

	q.Append("TEST QUERY")

	if q.String() != "TEST QUERY\n" {
		t.Errorf("Expected 'TEST QUERY\\n', got %q", q.String())
	}
}

func TestComplexQuery(t *testing.T) {
	q := createTestQuery()

	// Build a more complex query
	q.Match("n", "Person", "id", "123")
	q.With("n")
	q.Match("m", "Person", "id", "456")
	q.Relate("n", "KNOWS", "m")
	q.Set("n", map[string]string{"lastSeen": "2023-01-01"})
	q.Return()

	// Check that the generated query looks correct (ignoring exact spacing details)
	result := q.String()
	expected := [...]string{
		"MATCH (n:Person { id: '123' })",
		"WITH n",
		"MATCH (m:Person { id: '456' })",
		"MERGE (n)-[:KNOWS]->(m)",
		"SET n.lastSeen = '2023-01-01'",
		"RETURN count(*)",
	}

	for _, part := range expected {
		if !strings.Contains(result, part) {
			t.Errorf("Expected query to contain %q, but it doesn't: %q", part, result)
		}
	}
}

func TestCypherCharacterEscaping(t *testing.T) {
	q := createTestQuery()

	// Test a path with backslashes that needs to be fixed
	q.Match("n", "File", "path", `C:\Windows\System32\cmd.exe`)

	expected := "MATCH (n:File { path: 'c:/windows/system32/cmd.exe' }) "
	if q.String() != expected {
		t.Errorf("Expected query %q, got %q", expected, q.String())
	}
}

func TestNewQueryWithoutDriver(t *testing.T) {
	// Temporarily set Driver to nil
	originalDriver := Driver
	Driver = nil
	defer func() { Driver = originalDriver }()

	// This should return an error since driver is nil
	_, err := NewQuery()

	if err == nil {
		t.Error("Expected error when creating query without initialized driver, got nil")
	}
}

// TestInitDriverMock tests that we can mock the InitDriver function
// Note: This doesn't actually connect to a database
func TestInitDriverMockable(t *testing.T) {
	t.Skip("Skipping test that would require a Neo4j connection")
	/*
		err := InitDriver("bolt://localhost:7687", "neo4j", "password")
		if err != nil {
			t.Errorf("Failed to initialize driver: %v", err)
		}

		if Driver == nil {
			t.Error("Driver should be initialized after InitDriver")
		}
	*/
}
