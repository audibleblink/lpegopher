package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usage: gump.exe <service | task>")
	}

	switch os.Args[1] {
	case "task":
		listTasks()
	case "service":
		listServices()
	}
}

func evaluatePath(path string) (out string) {
	// https://gitlab.com/stu0292/windowspathenv
	out = path
	// return unless strings starts with %
	if !strings.HasPrefix(path, "%") {
		return
	}
	// return unless there's a second %
	trim := strings.TrimPrefix(path, "%")
	i := strings.Index(trim, "%")
	if i == -1 {
		return
	}
	// check if substr between two % is the name of an existing env var
	val, ok := os.LookupEnv(trim[:i])
	if !ok {
		return
	}
	// env var value will use os path separator
	remainder := filepath.FromSlash(trim[i+1:])
	// check the remainder starts with path separateor
	if !strings.HasPrefix(remainder, "\\") {
		return
	}
	// prepend the value to the remainder of the path
	return val + remainder
}
