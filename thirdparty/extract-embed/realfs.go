package extractembed

import (
	"os"
	"strings"
)

func saveFile(name string, body []byte, base string) error {
	makeSuperFolder(name, base)

	f, err := os.Create(base + "/" + name)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write(body)

	if err != nil {
		return err
	}

	return nil
}

func makeSuperFolder(path string, base string) error {
	splited := strings.Split(base+"/"+path, "/")
	splited[len(splited)-1] = ""

	joined := strings.Join(splited, "/")

	err := os.MkdirAll(joined, 0777)

	return err
}
