package collectors

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Microsoft/go-winio"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
	"github.com/audibleblink/pegopher/util"
	winacl "github.com/kgoins/go-winacl/pkg"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/windows"
	"www.velocidex.com/golang/binparsergen/reader"
	"www.velocidex.com/golang/go-pe"
)

type DACL struct {
	Owner string        `json:"Owner"`
	Group string        `json:"Group"`
	Aces  []ReadableAce `json:"Aces"`
}

// INode contains the parsed import and exports of the INode
type INode struct {
	Name     string   `json:"Name"`
	Path     string   `json:"Path"`
	Parent   string   `json:"Dir"`
	Type     string   `json:"Type"`
	Imports  []string `json:"Imports"`
	Forwards []string `json:"Forwards"`
	DACL     DACL     `json:"DACL"`
}

type ReadableAce struct {
	Principal string   `json:"Principal"`
	Rights    []string `json:"Rights"`
}

func PEs(writer io.Writer, dir string) {
	log := logerr.Add("pe collector")

	walkStartPath, _ := filepath.Abs(dir)
	walkFunction := walkFunctionGenerator(writer)

	objs, err := os.ReadDir(walkStartPath)
	if err != nil {
		log.Fatalf("could not read %s", walkStartPath)
	}

	g := new(errgroup.Group)
	for _, obj := range objs {
		if obj.IsDir() {
			path := filepath.Join(walkStartPath, obj.Name())
			log.Debugf("forking collection of %s", path)
			g.Go(func() error {
				err := filepath.WalkDir(path, walkFunction)
				log.Infof("completed collection of %s", path)
				return err
			})
		}
	}
	if err := g.Wait(); err == nil {
		logerr.Info("All collection jobs complete")
	}
}

func walkFunctionGenerator(writer io.Writer) fs.WalkDirFunc {
	// use a set to track if a report for a PE's parent directory
	// has already been printed
	printedParentDir := &sync.Map{}

	return func(path string, info os.DirEntry, err error) error {
		log := logerr.Add("dirwalk")

		if err != nil {
			log.Warnf("recursion, amirite?: %s", err)
		}

		if info.IsDir() {
			return nil
		}

		path = util.Lower(path)
		isExe, _ := filepath.Match("*.exe", filepath.Base(path))
		isDll, _ := filepath.Match("*.dll", filepath.Base(path))

		if isExe || isDll {
			parent := filepath.Dir(path)
			_, alreadyDidIt := printedParentDir.LoadOrStore(parent, true)
			if !alreadyDidIt {
				dirReport := newDirectoryReport(parent)
				jsPrint(writer, dirReport)
			}

			report := newPEReport(path)
			report.Parent = parent

			peFile, err := newPEFile(report.Path)
			if err != nil {
				log.Debugf("pe parsing failed: %s", err)
				return nil
			}

			err = populatePEReport(report, peFile)
			if err != nil {
				log.Warnf("could not generate report for %s: %s", path, err)
				return nil
			}

			jsPrint(writer, report)

		}
		return nil
	}
}

func newDirectoryReport(path string) *INode {
	report := &INode{}
	report.Name = filepath.Base(path)
	report.Path, _ = filepath.Abs(path)
	report.Type = node.Dir
	report.Parent = filepath.Dir(path)
	err := handleDirPerms(report)
	if err != nil {
		return report
	}
	return report
}

func newPEReport(path string) *INode {
	report := &INode{}
	report.Name = filepath.Base(path)
	report.Path, _ = filepath.Abs(path)

	if strings.HasSuffix(util.Lower(report.Path), ".dll") {
		report.Type = node.Dll
	} else if strings.HasSuffix(util.Lower(report.Path), ".exe") {
		report.Type = node.Exe
	}
	return report
}

func newPEFile(path string) (pefile *pe.PEFile, err error) {
	peFileH, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return
	}

	peReader, err := reader.NewPagedReader(peFileH, 4096, 100)
	if err != nil {
		return
	}

	pefile, err = pe.NewPEFile(peReader)
	return
}

func populatePEReport(report *INode, peFile *pe.PEFile) error {
	report.Imports = peFile.Imports()
	report.Forwards = peFile.Forwards()

	dacl, err := pullDACL(report.Path)
	if err != nil {
		return err
	}
	report.DACL = dacl
	return nil
}

func pullDACL(path string) (DACL, error) {
	dacl := DACL{}
	sd, err := securityDescriptorFor(path)
	if err != nil {
		return dacl, err
	}
	dacl.Owner = sidResolve(sd.Owner)
	dacl.Group = sidResolve(sd.Group)
	for _, ace := range sd.DACL.Aces {
		dacl.Aces = append(dacl.Aces, newReadableAce(ace))
	}
	return dacl, err
}

func securityDescriptorFor(path string) (sd winacl.NtSecurityDescriptor, err error) {
	winSD, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if !winSD.IsValid() {
		return sd, fmt.Errorf("invalid security descriptor %s", err)
	}

	// convert windows.SD into SDDL, then back into an SD
	// represented as a byte slice, so go-winacl can parse it
	sdBytes, err := winio.SddlToSecurityDescriptor(winSD.String())
	if err != nil {
		return
	}

	sd, err = winacl.NewNtSecurityDescriptor(sdBytes)
	return
}

func newReadableAce(ace winacl.ACE) ReadableAce {
	var rAce ReadableAce

	perms := ace.AccessMask.String()
	rAce.Rights = strings.Split(perms, " ")

	switch ace.ObjectAce.(type) {
	case winacl.BasicAce:
		rAce.Principal = sidResolve(ace.ObjectAce.GetPrincipal())

	case winacl.AdvancedAce:
		aa := ace.ObjectAce.(winacl.AdvancedAce)
		sid := aa.GetPrincipal()
		rAce.Principal = sidResolve(sid)
	}
	return rAce
}

func sidResolve(sid winacl.SID) string {
	res := sid.Resolve()
	if strings.HasPrefix(res, "S-1-") {
		// failed to resolve
		winSID, err := windows.StringToSid(sid.String())
		if err != nil {
			return res
		}
		user, domain, _, err := winSID.LookupAccount("")
		if err != nil {
			return res
		}
		return fmt.Sprintf(`%s\%s`, domain, user)
	}
	return res
}

func handleDirPerms(report *INode) error {
	dacl, err := pullDACL(report.Path)
	if err != nil {
		return err
	}
	report.DACL = dacl
	return nil
}

func jsPrint(writer io.Writer, report *INode) {
	serialized, _ := json.Marshal(report)
	fmt.Fprintln(writer, string(serialized))
}
