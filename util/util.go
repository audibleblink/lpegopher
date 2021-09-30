package util

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
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

func DecodeUTF16(b []byte) ([]byte, error) {

	if bytes.HasPrefix(b, []byte{0xff, 0xfe}) {
		b = b[2:]
	}

	if bytes.HasPrefix(b, []byte{0x00}) {
		b = b[1:]
	}

	if len(b)%2 != 0 {
		return []byte{}, fmt.Errorf("must have even length byte slice")
	}

	u16s := make([]uint16, 1)
	buffer := &bytes.Buffer{}
	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		buffer.Write(b8buf[:n])
	}

	newBytes := buffer.Bytes()

	return newBytes, nil
}
