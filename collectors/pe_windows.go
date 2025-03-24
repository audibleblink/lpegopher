package collectors

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Microsoft/go-winio"
	"github.com/audibleblink/go-winacl"
	"golang.org/x/sys/windows"
	"www.velocidex.com/golang/binparsergen/reader"
	"www.velocidex.com/golang/go-pe"

	"github.com/audibleblink/lpegopher/logerr"
	"github.com/audibleblink/lpegopher/node"
	"github.com/audibleblink/lpegopher/util"
)

func PEs(dir string) {
	log := logerr.Add("pe collector")
	walkStartPath, _ := filepath.Abs(dir)
	filepath.WalkDir(walkStartPath, walkFunction)
	log.Infof("completed collection of %s", walkStartPath)
}

func walkFunction(path string, info os.DirEntry, err error) error {
	log := logerr.Add("dirwalk")

	if err != nil {
		log.Warnf("%v", err)
	}

	if info.IsDir() {
		return nil
	}

	path = util.Lower(path)
	isExe, _ := filepath.Match("*.exe", filepath.Base(path))
	isDll, _ := filepath.Match("*.dll", filepath.Base(path))

	if isExe || isDll {
		parent := filepath.Dir(path)
		_, alreadyDidIt := cache.LoadOrStore(parent, true)
		if !alreadyDidIt {
			dirReport := newDirectoryReport(parent)
			doPrint(dirReport)
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

		doPrint(report)
	}
	return nil
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
	imports := make([]*Dep, 0)
	for _, imp := range peFile.Imports() {
		imports = append(imports, &Dep{Name: util.Lower(imp)})
	}
	report.Imports = imports

	forwards := make([]*Dep, 0)
	for _, fwd := range peFile.Forwards() {
		forwards = append(forwards, &Dep{Name: util.Lower(fwd)})
	}
	report.Forwards = forwards

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
	dacl.Owner = &Principal{Name: sidResolve(sd.Owner)}
	dacl.Group = &Principal{Name: sidResolve(sd.Group)}
	for _, ace := range sd.DACL.Aces {
		dacl.Aces = append(dacl.Aces, newReadableAce(ace))
	}
	return dacl, err
}

func securityDescriptorFor(path string) (sd winacl.NtSecurityDescriptor, err error) {
	winSD, err := windows.GetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION,
	)
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
		name := sidResolve(ace.ObjectAce.GetPrincipal())
		rAce.Principal = &Principal{Name: name}

	case winacl.AdvancedAce:
		aa := ace.ObjectAce.(winacl.AdvancedAce)
		sid := aa.GetPrincipal()
		name := sidResolve(sid)
		rAce.Principal = &Principal{Name: name}
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

func doPrint(report *INode) {
	var nodeID string
	switch report.Type {
	case node.Exe:
		nodeID = report.Write(writers[ExeFile])
	case node.Dll:
		nodeID = report.Write(writers[DllFile])
	case node.Dir:
		nodeID = report.Write(writers[DirFile])
	}

	for _, ace := range report.DACL.Aces {
		pID := ace.Principal.Write(writers[PrincipalFile])
		for _, priv := range ace.Rights {
			if node.AbusableAces[priv] {
				rel := &Rel{
					Start: pID,
					Rel:   priv,
					End:   nodeID,
				}
				rel.Write(writers[RelsFile])
			}
		}
	}

	for _, fwd := range report.Forwards {
		re := regexp.MustCompile(`\..*$`)
		fwd.Name = re.ReplaceAllLiteralString(fwd.Name, ".dll")
		fwdID := fwd.Write(writers[DepsFile])
		rel := &Rel{
			Start: nodeID,
			Rel:   Forwards,
			End:   fwdID,
		}
		rel.Write(writers[RelsFile])
	}

	for _, imp := range report.Imports {
		re := regexp.MustCompile(`!.*$`)
		imp.Name = re.ReplaceAllLiteralString(imp.Name, "")
		impID := imp.Write(writers[DepsFile])
		rel := &Rel{
			Start: nodeID,
			Rel:   Imports,
			End:   impID,
		}
		rel.Write(writers[ImportFile])
	}
}
