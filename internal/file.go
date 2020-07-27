package internal

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
)

// AllFilesInDirectory return all file in directory recursive
func AllFilesInDirectory(directory string) ([]string, error) {
	var files []string
	return allFilesInDirectory(directory, files)
}

func allFilesInDirectory(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = allFilesInDirectory(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

// CreateZipArchive create a zip archive file from files
func CreateZipArchive(files map[string][]byte) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)
	defer func() {
		_ = zipWriter.Close()
	}()

	for name, content := range files {
		writer, err := zipWriter.CreateHeader(&zip.FileHeader{Name: name})
		if err != nil {
			return nil, err
		}

		_, err = writer.Write(content)
		if err != nil {
			return nil, err
		}
	}

	return &buffer, nil
}