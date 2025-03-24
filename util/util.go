package util

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(1337)

// Lower converts a string to lowercase
func Lower(str string) string {
	return strings.ToLower(str)
}

// PathFix sanitizes path strings for use in the database
func PathFix(str string) string {
	str = strings.ReplaceAll(str, `"`, "")
	str = resolveEnvPath(str)
	str = strings.ReplaceAll(str, `\`, "/")
	str = strings.ReplaceAll(str, `,`, `.`)
	str = strings.Trim(str, " ")
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

// EvaluatePath resolves environment variables in a path
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

// LineCount counts lines in a reader
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

// Rand generates a random string
func Rand() string {
	size := 6
	sb := strings.Builder{}
	sb.Grow(size)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := size-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// SmoothBrainPath splits a command line into executable and arguments
func SmoothBrainPath(cmdline string) (bin, args string) {
	if strings.HasPrefix(cmdline, `"`) {
		quoteCharOffset := 1
		secondQuoteIdx := strings.Index(cmdline[quoteCharOffset:], `"`)
		endOfCmd := len(cmdline) - len(cmdline[secondQuoteIdx:]) + quoteCharOffset
		bin = cmdline[quoteCharOffset:endOfCmd]
		args = strings.TrimSpace(cmdline[endOfCmd+quoteCharOffset:])
		return
	}

	splitCmd := strings.Split(cmdline, " ")

	for idx, part := range splitCmd {
		if strings.HasSuffix(part, ".exe") || strings.HasSuffix(part, ".dll") {
			bin = strings.Join(splitCmd[0:idx+1], " ")
			args = strings.Join(splitCmd[idx+1:], " ")
		}
	}

	return
}
