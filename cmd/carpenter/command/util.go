package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func getErrorMessages(errs []error) []string {
	msg := make([]string, 0, len(errs))
	for _, e := range errs {
		msg = append(msg, e.Error())
	}
	return msg
}

func execute(queries []string) error {
	for _, query := range queries {
		if !dryrun {
			if _, err := db.Exec(query); err != nil {
				return fmt.Errorf("err: db.Exec `%s' failed for reason %s", query, err)
			}
		}
		if verbose {
			fmt.Println(query)
		}
	}
	return nil
}

func walk(path, ext string) (map[string][]string, error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("err: ioutil.ReadDir %s for reason %v", path, err)
	}
	files := []os.FileInfo{}
	for _, file := range dir {
		filename := file.Name()
		pos := strings.LastIndex(filename, ".")
		if pos <= 0 {
			continue
		}
		if filename[pos:] != ext {
			continue
		}
		files = append(files, file)
	}
	if len(files) <= 0 {
		return nil, fmt.Errorf("err: No csv files found in %s", path)
	}

	filesMap := map[string][]string{}
	for _, file := range files {
		filename := file.Name()
		splited := strings.Split(filename, string(os.PathSeparator))
		table := strings.Split(splited[len(splited)-1], ".")[0]
		if _, ok := filesMap[table]; !ok {
			filesMap[table] = []string{}
		}
		filesMap[table] = append(filesMap[table], fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), filename))
	}
	return filesMap, nil
}
