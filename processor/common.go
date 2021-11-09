package processor

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

var (
	MaxBatchSize = 1
)

var (
	CurrentBatchLen = 1
	CurrentBatch    = &strings.Builder{}
	BatchCounter    = 0
)

type queryBuilder func([]byte) (*cypher.Query, error)

func QueryBuilder(callback queryBuilder) func(string) error {

	return func(path string) error {
		log := logerr.Add("querybuilder")
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		lineCount, err := util.LineCount(file)
		if err != nil {
			log.Fatalf("could not get line count for %s", path)
		}
		log.Infof("processing %d lines with batch size %d == %d batches", lineCount, MaxBatchSize, lineCount*1.0/MaxBatchSize+1)

		file.Seek(0, 0)

		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 8*1024)
		scanner.Buffer(buf, 1024*1024)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber += 1
			text := scanner.Bytes()

			cypherQ, err := callback(text)
			if err != nil {
				log.Errorf("error generating query for line %d", lineNumber)
				continue
			}

			// CurrentBatch.WriteString(cypherQ.String())
			if CurrentBatchLen%MaxBatchSize == 0 || lineNumber == lineCount {
				BatchCounter++
				log.Infof("commiting batch transaction %d", BatchCounter)
				CurrentBatch.WriteString(cypherQ.String())
				cypherQ.Raw(CurrentBatch.String()).Return()
				err = cypherQ.ExecuteW()
				if err != nil {
					err = errors.Unwrap(err)
					switch e := err.(type) {
					case *neo4j.Neo4jError:
						if e.Code == "Neo.ClientError.Schema.ConstraintValidationFailed" {
							log.Debugf("duplicate entry (%s)", e.Msg)
							continue
						} else {
							log.Errorf("neo  %s", err)
							// continue
						}
					default:
						log.Errorf("query failed %s", err)
						continue
					}
				}
				CurrentBatch.Reset()
				CurrentBatchLen = 1
			} else {
				if len(cypherQ.String()) > 0 {
					CurrentBatch.WriteString(cypherQ.String())
					CurrentBatchLen++
				}
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
