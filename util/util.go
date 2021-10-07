package util

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Lower(str string) string {
	return strings.ToLower(str)
}

func PathFix(str string) string {
	str = strings.Trim(str, `"`)
	str = resolveEnvPath(str)
	str = strings.ReplaceAll(str, `\`, "/")
	// swap slack direction to avoid cross-platform issues
	return Lower(str)
}

func resolveEnvPath(path string) (out string) {

	// return the original filepath unchanged unless we get to the end
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

// EvaluatePath will resolve any environment variables in a path string
func EvaluatePath(path string) (out string) {
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

func LineCount(r io.Reader) (int, error) {
	buffer := make([]byte, 32*1024)
	lineSep := []byte{'\n'}
	count := 0

	for {

		byteCount, err := r.Read(buffer)
		count += bytes.Count(buffer[:byteCount], lineSep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}

	}
}
