package processor

import (
	"bufio"
	"os"
	"strings"

	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
)

const (
	MaxBatchSize = 5000
)

var (
	BatchCounter = 1
	CurrentBatch = &strings.Builder{}
)

type queryBuilder func([]byte) (*cypher.Query, error)

func QueryBuilder(callback queryBuilder) func(string) error {

	return func(path string) error {
		log := logerr.Add("relate generator")
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		lineCount, err := util.LineCount(file)
		if err != nil {
			log.Fatalf("could not get line count for %s", path)
		}

		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 8*1024)
		scanner.Buffer(buf, 1024*1024)
		count := 0
		for scanner.Scan() {
			count += 1
			text := scanner.Bytes()

			cypherQ, err := callback(text)
			if err != nil {
				log.Errorf("error generating 	query for line %d", count)
				continue
			}

			if BatchCounter%MaxBatchSize == 0 || count == lineCount {
				cypherQ.Raw(CurrentBatch.String())
				CurrentBatch.Reset()
				BatchCounter = 1
				err = cypherQ.ExecuteW()
				if err != nil {
					log.Errorf("query failed %b", err)
				}
			} else {
				CurrentBatch.WriteString(cypherQ.Terminate().String())
				BatchCounter++
			}
		}

		if err := scanner.Err(); err != nil {
			log.Errorf("scanner quit: %s", err)
		}
		return err
	}

}

//
// func NewFileProcessor(callback queryBuilder) func(file string) error {
// 	return func(path string) error {
// 		file, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}
//
// 		count := 0
// 		scanner := bufio.NewScanner(file)
// 		buf := make([]byte, 0, 8*1024)
// 		scanner.Buffer(buf, 1024*1024)
// 		for scanner.Scan() {
// 			count += 1
// 			text := scanner.Bytes()
// 			cypherQ, err := callback(text)
// 			if err != nil {
// 				switch err.(type) {
// 				case *json.SyntaxError:
// 					fmt.Fprintf(os.Stderr, "malformed json at line %d", count)
// 					continue
// 				default:
// 					return err
// 				}
// 			}
//
// 		}
//
// 		if err := scanner.Err(); err != nil {
// 			logerr.Errorf("scanner quit: %s", err)
// 		}
//
// 		logerr.Infof("Done. Processed %d lines", count)
// 		return nil
//
// 	}
// }
