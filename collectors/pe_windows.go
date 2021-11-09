package collectors

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Microsoft/go-winio"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/node"
	"github.com/audibleblink/pegopher/util"
	winacl "github.com/kgoins/go-winacl/pkg"
	"golang.org/x/sys/windows"
	"www.velocidex.com/golang/binparsergen/reader"
	"www.velocidex.com/golang/go-pe"
)

type PEFunction struct {
	Host      string   `json:"Host"`
	Functions []string `json:"Functions"`
}

type DACL struct {
	Owner string        `json:"Owner"`
	Group string        `json:"Group"`
	Aces  []ReadableAce `json:"Aces"`
}

// INode contains the parsed import and exports of the INode
type INode struct {
	Name     string       `json:"Name"`
	Path     string       `json:"Path"`
	Parent   string       `json:"Dir"`
	Type     string       `json:"Type"`
	ImpHash  string       `json:"ImpHash"`
	Exports  []string     `json:"Exports"`
	Imports  []PEFunction `json:"Imports"`
	Forwards []PEFunction `json:"Forwards"`
	DACL     DACL         `json:"DACL"`
}

type ReadableAce struct {
	Principal string   `json:"Principal"`
	Rights    []string `json:"Rights"`
}

func PEs(writer io.Writer, dir string) {
	walkStartPath, _ := filepath.Abs(dir)
	walkFunction := walkFunctionGenerator(writer)
	filepath.WalkDir(walkStartPath, walkFunction)
}

func walkFunctionGenerator(writer io.Writer) fs.WalkDirFunc {

	return func(path string, info os.DirEntry, err error) error {
		// use a set to track if a report for a PE's parent directory
		// has already been printed
		printedParentDir := make(map[string]bool)
		log := logerr.Add("dirwalk")

		if err != nil {
			log.Warnf("HUH", err)
		}

		if info.IsDir() {
			return nil
		}

		path = util.Lower(path)
		isExe, _ := filepath.Match("*.exe", filepath.Base(path))
		isDll, _ := filepath.Match("*.dll", filepath.Base(path))

		if isExe || isDll {
			parent := filepath.Dir(path)
			if !printedParentDir[parent] {
				// first time finding a PE in this directory
				dirReport := newDirectoryReport(parent)
				jsPrint(writer, dirReport)
				printedParentDir[parent] = true
			}

			report := newPEReport(path)
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

	return pe.NewPEFile(peReader)
}

func populatePEReport(report *INode, peFile *pe.PEFile) error {
	report.ImpHash = peFile.ImpHash()
	report.Imports = genPEFunctions(peFile.Imports())
	report.Forwards = genPEFunctions(patchForwards(peFile.Forwards()))
	report.Exports = patchExports(peFile.Exports())
	report.Parent = filepath.Dir(report.Path)

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

func patchExports(funcs []string) (out []string) {
	for _, fun := range funcs {
		// strip leading ':'
		out = append(out, fun[1:])
	}
	return
}

func patchForwards(funcs []string) (out []string) {
	for _, fun := range funcs {
		// dbgcore.MiniDumpWriteDump....
		matcher := regexp.MustCompile(`\.`)
		s := matcher.ReplaceAllString(fun, ".dll!")
		out = append(out, s)
	}
	return
}
func genPEFunctions(list []string) []PEFunction {
	// incoming: ["dllname!funcName"]
	funcs := []PEFunction{}
	accumulatedFns := make(map[string][]string)
	for _, fn := range list {
		splitFn := strings.Split(fn, "!")
		peName := splitFn[0]
		funcName := splitFn[1]
		accumulatedFns[peName] = append(accumulatedFns[peName], funcName)
	}

	for peName, funcSlice := range accumulatedFns {
		funcs = append(funcs, PEFunction{peName, funcSlice})
	}
	return funcs
}

func jsPrint(writer io.Writer, report *INode) {
	serialized, _ := json.Marshal(report)
	fmt.Fprintln(writer, string(serialized))
}
