package internal

import (
	"fmt"
	"github.com/mylxsw/go-toolkit/container"
	"github.com/mylxsw/go-toolkit/log"
	"io/ioutil"
	"os"
	"strings"
)

var logger = log.Module("artisan")

func Artisan(cc *container.Container, skeleton, dest string, data Data) error {
	files, err := AllFilesInDirectory("./skeleton")
	if err != nil {
		return fmt.Errorf("read .sk files from %s failed: %s", skeleton, err)
	}

	parsedFiles := make(map[string]string)
	for _, f := range files {
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
