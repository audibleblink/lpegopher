package processor

import (
	"bufio"
	"errors"
	"os"

	"github.com/audibleblink/pegopher/cypher"
	"github.com/audibleblink/pegopher/logerr"
	"github.com/audibleblink/pegopher/util"
)

const goRoutineLimit = 50

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
		log.Infof("processing %d inode entries", lineCount)
		file.Seek(0, 0)

		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 2048*2048)
		scanner.Buffer(buf, len(buf))
		lineNumber := 0
		checkpoint := lineCount / 20

		wg := util.NewLimitedWaitGroup(goRoutineLimit)
		for scanner.Scan() {
			lineNumber += 1
			if lineNumber%checkpoint == 0 {
				log.Infof("%d lines processed - %f%% done...", lineNumber, (float32(lineNumber)/float32(lineCount))*100)
			}

			text := scanner.Bytes()

			cypherQ, err := callback(text)
			if err != nil {
				log.Errorf("error generating query for line %d", lineNumber)
				continue
			}

			wg.Add(1)
			go (func(query *cypher.Query, lwg *util.LimitedWaitGroup) error {
				err = query.Return().ExecuteW()
				if err != nil {
					err = errors.Unwrap(err)
					log.Errorf("query failed %s", err)
				}
				lwg.Done()
				return err
			})(cypherQ, wg)
		}

		if err := scanner.Err(); err != nil {
			log.Errorf("scanner quit on line %d: %s", lineNumber, err)
		}

		log.Info("waiting on neo4j to report completion")
		wg.Wait()
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
