package extractembed

import (
	"embed"
	"fmt"
)

func listEmbedFiles(path string, embedFs *embed.FS) ([]string, error) {
	fileList := []string{}

	dir, err := embedFs.ReadDir(path)

	if err != nil {
		return []string{}, fmt.Errorf("error reading dir: %w", err)
	}

	if path == "." {
		path = ""
	} else {
		path = path + "/"
	}

	for _, file := range dir {
		filePath := path + file.Name()

		if file.IsDir() {
			fil, err := listEmbedFiles(filePath, embedFs)

			if err != nil {
				return []string{}, err
			}

			fileList = append(fileList, fil...)

		} else {
			fileList = append(fileList, filePath)
		}
	}

	return fileList, nil
}

func readEmbedFile(filePath string, embedFs *embed.FS) ([]byte, error) {
	body, err := embedFs.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}
