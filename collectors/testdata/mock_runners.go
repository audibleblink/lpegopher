package testdata

import (
	"github.com/audibleblink/lpegopher/collectors"
)

// MockRunnerGenerator creates test PERunner instances
func MockRunnerGenerator() []collectors.PERunner {
	// Mock Principal for the context
	systemPrincipal := &collectors.Principal{
		Name: "SYSTEM",
		Type: "user",
	}

	// Mock INode for the executable
	testExe := &collectors.INode{
		Name:   "test.exe",
		Path:   "C:/Windows/System32/test.exe",
		Parent: "C:/Windows/System32",
		Type:   "exe",
	}

	// Create a test service
	testService := collectors.PERunner{
		Name:     "TestService",
		Type:     "service",
		Args:     "-k testService",
		RunLevel: "system",
		Exe:      testExe,
		Context:  systemPrincipal,
	}

	// Create a test scheduled task
	testTask := collectors.PERunner{
		Name:     "TestTask",
		Type:     "task",
		Args:     "/param1 /param2",
		RunLevel: "user",
		Exe:      testExe,
		Context:  systemPrincipal,
	}

	// Create a test autorun entry
	testAutorun := collectors.PERunner{
		Name:     "TestAutorun",
		Type:     "autorun",
		Args:     "/autostart",
		RunLevel: "user",
		Exe:      testExe,
		Context:  systemPrincipal,
	}

	return []collectors.PERunner{testService, testTask, testAutorun}
}
