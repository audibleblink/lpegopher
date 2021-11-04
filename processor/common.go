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
	MaxBatchSize = 5
)

var (
	CurrentBatchLen = 1
	CurrentBatch    = &strings.Builder{}
	BatchCounter    = 1
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
		count := 0
		for scanner.Scan() {
			count += 1
			text := scanner.Bytes()

			cypherQ, err := callback(text)
			if err != nil {
				log.Errorf("error generating query for line %d", count)
				continue
			}
			if CurrentBatchLen%MaxBatchSize == 0 || count == lineCount {
				log.Infof("commiting batch transaction %d", BatchCounter)
				cypherQ.Raw(CurrentBatch.String())
				cypherQ.Return()
				err = cypherQ.ExecuteW()
				if err != nil {
					err = errors.Unwrap(err)
					switch e := err.(type) {
					case *neo4j.Neo4jError:
						if e.Code == "Neo.ClientError.Schema.ConstraintValidationFailed" {
							log.Debugf("duplicate entry (%s)", e.Msg)
						}
					default:
						log.Errorf("query failed %s", err)
					}
					log.Errorf("query failed %s", err)
					continue
				}
				BatchCounter++
				CurrentBatch.Reset()
				CurrentBatchLen = 1
			} else {
				if len(cypherQ.String()) > 0 {
					CurrentBatch.WriteString(cypherQ.Terminate().String())
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
