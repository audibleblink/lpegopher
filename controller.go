package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/alexflint/go-arg"
)

type jsonProcessor func([]byte) error

func handleProcess(cli *arg.Parser, args argType) (err error) {
	switch {
	case args.Process.Dlls != nil:
	case args.Process.Exes != nil:
	case args.Process.Tasks != nil:
		processor := newFileProcessor(NewRunnerFromJson)
		err = processor(args.Process.Tasks.File)
	case args.Process.Services != nil:
		processor := newFileProcessor(NewRunnerFromJson)
		err = processor(args.Process.Services.File)
	}
	return
}

func processDlls(args argType)  {}
func processExes(args argType)  {}
func processTasks(args argType) {}

func newFileProcessor(jp jsonProcessor) func(file string) error {
	return func(path string) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Bytes()
			text, err = DecodeUTF16(text)
			if err != nil {
				return err
			}
			err = jp(text)
			if err != nil {
				return err
			}
		}
		return nil

	}
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
