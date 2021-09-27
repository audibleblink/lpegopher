package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	securableIFile

	ContainedExes []*EXE  `gogm:"direction=outgoing;relationship=CONTAINS"`
	ContainedDlls []*DLL  `gogm:"direction=outgoing;relationship=CONTAINS"`
	HostsPEs      *Runner `gogm:"direction=outgoing;relationship=HOSTS_PES"`
}

func (d *Directory) HostsPEsFor(runner *Runner) error {
	d.HostsPEs = runner
	return d.save()
}

func (d *Directory) Add(ifile interface{}) error {
	switch f := ifile.(type) {
	case *EXE:
		d.ContainedExes = append(d.ContainedExes, f)
	case *DLL:
		d.ContainedDlls = append(d.ContainedDlls, f)
	}

	return d.save()
}

// Merge will either create or retreive the node based on the key-valie pair provides
// In this case, the Directory struct designates the "path" field as unique
func (x *Directory) Merge(uniquePropName, propValue string) (err error) {
	nodeType := "Directory"
	sess, err := newNeoSession()
	if err != nil {
		return err
	}

	queryTemplate := `MERGE (x:%s {%s: "%s"}) RETURN x`
	query := fmt.Sprintf(queryTemplate, nodeType, uniquePropName, (propValue))
	err = sess.Query(context.Background(), query, nil, x)
	return sess.Save(context.Background(), x)
}

func (x *Directory) SetName(name string) error {
	x.Name = name
	return x.save()
}

func (x *Directory) save() (err error) {
	if x.Id == nil {
		return fmt.Errorf("no ID provided. ensure this node exists before attempting to update a property")
	}
	sess, err := newNeoSession()
	if err != nil {
		return err
	}
	return sess.Save(context.Background(), x)
}

func lower(str string) string {
	return strings.ToLower(str)
}

func pathFix(str string) string {
	str = strings.Trim(str, `"`)
	str = resolveEnvPath(str)
	str = strings.ReplaceAll(str, "\\", "/")
	// swap slack direction to avoid cross-platform issues
	return lower(str)
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
