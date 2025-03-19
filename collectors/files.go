package collectors

import (
	"encoding/hex"
	"os"
	"sync"

	"github.com/audibleblink/concurrent-writer"
	"github.com/audibleblink/lpegopher/logerr"
)

// Constants for file paths used for outputs
const (
	ExeFile       = "exes.csv"          // Path to write executable file data
	DllFile       = "dlls.csv"          // Path to write dynamic link library data
	DirFile       = "dirs.csv"          // Path to write directory data
	PrincipalFile = "principals.csv"    // Path to write security principal data
	RelsFile      = "relationships.csv" // Path to write relationship data
	DepsFile      = "deps.csv"          // Path to write dependency data
	RunnersFile   = "runners.csv"       // Path to write auto-runner data
	ImportFile    = "imports.csv"       // Path to write import relationship data
)

var (
	key, _ = hex.DecodeString("900F02030405060708090A0B9C0D0E0FF0E0D0C0B0A090807060504030201091")
	cache  = &sync.Map{}
)

var (
	writers                        map[string]*concurrent.Writer
	f0, f1, f2, f3, f4, f5, f6, f7 os.File
)

// InitOutputFiles initializes output files for data collection
func InitOutputFiles() {

	var (
		f0, _ = os.OpenFile(ExeFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		f1, _ = os.OpenFile(DllFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		f2, _ = os.OpenFile(DirFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		f3, _ = os.OpenFile(PrincipalFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		f4, _ = os.OpenFile(RelsFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		f5, _ = os.OpenFile(DepsFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		f6, _ = os.OpenFile(RunnersFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		f7, _ = os.OpenFile(ImportFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	)

	writers = map[string]*concurrent.Writer{
		ExeFile:       concurrent.NewWriter(f0),
		DllFile:       concurrent.NewWriter(f1),
		DirFile:       concurrent.NewWriter(f2),
		PrincipalFile: concurrent.NewWriter(f3),
		RelsFile:      concurrent.NewWriter(f4),
		DepsFile:      concurrent.NewWriter(f5),
		RunnersFile:   concurrent.NewWriter(f6),
		ImportFile:    concurrent.NewWriter(f7),
	}
}

// FlashAndClose flushes all writer buffers and closes files
func FlashAndClose() {
	log := logerr.Add("cleanup")

	defer f0.Close()
	defer f1.Close()
	defer f2.Close()
	defer f3.Close()
	defer f4.Close()
	defer f5.Close()
	defer f6.Close()
	defer f7.Close()

	for f, writer := range writers {
		err := writer.Flush()
		if err != nil {
			log.Errorf("could not flush %s: %v", f, err)
			continue
		}
	}

}
