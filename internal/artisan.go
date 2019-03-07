package internal

import (
	"fmt"
	"github.com/mylxsw/go-toolkit/collection"
	"github.com/mylxsw/go-toolkit/container"
	"github.com/mylxsw/go-toolkit/log"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var logger = log.Module("artisan")

func Artisan(cc *container.Container, skeleton, dest string, data Data) error {
	files, err := AllFilesInDirectory(skeleton)
	if err != nil {
		return fmt.Errorf("read .sk files from %s failed: %s", skeleton, err)
	}

	excludeFiles := make([]string, 0)

	collection.MustNew(files).Filter(func(path string) bool {
		return strings.HasSuffix(path, "/exclude.skc")
	}).Each(func(ef string) {
		excludeBytes, err := ioutil.ReadFile(ef)
		if err != nil {
			panic(err)
		}

		lines, err := data.Parse(string(excludeBytes))
		if err != nil {
			panic(err)
		}

		collection.MustNew(strings.Split(lines, "\n")).Filter(func(line string) bool {
			return strings.TrimSpace(line) != ""
		}).Each(func(line string) {
			excludeFiles = append(excludeFiles, path.Join(path.Dir(ef), line))
		})
	})

	parsedFiles := make(map[string]string)
	for _, f := range files {
		if stringHasPrefixes(path.Join(path.Dir(f), path.Base(f)), excludeFiles) || strings.HasSuffix(f, "/exclude.skc") {
			continue
		}

		content, err := ioutil.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s failed: %s", f, err)
		}

		if strings.HasSuffix(f, ".sk") {
			res, err := data.Parse(string(content))
			if err != nil {
				return fmt.Errorf("parse file %s failed: %s", f, err)
			}

			parsedFiles[f[:len(f)-3]] = res
			logger.Debugf("parse %s -> %s ok", f, f[:len(f)-3])
		} else {
			parsedFiles[f] = string(content)
			logger.Debugf("copy %s ok", f)
		}
	}

	buffer, err := CreateZipArchive(parsedFiles)
	if err != nil {
		return fmt.Errorf("create zip archive failed: %s", err)
	}

	if err := ioutil.WriteFile(dest, buffer.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("save zip archive [%s] to filesystem failed: %s", dest, err)
	}

	logger.Infof("zip archive file has been saved to %s", dest)

	return nil
}

func stringHasPrefixes(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}

	return false
}
