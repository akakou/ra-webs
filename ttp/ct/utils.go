package ct

import (
	"io"
	"os"
)

func readFile(name string) (string, error) {
	lastFile, err := os.Open(name)
	if err != nil {
		return "", nil
	}

	defer lastFile.Close()

	b, err := io.ReadAll(lastFile)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func writeFile(body string) error {
	lastFile, err := os.Open(body)
	if err != nil {
		return err
	}

	defer lastFile.Close()

	_, err = io.WriteString(lastFile, body)
	if err != nil {
		return err
	}

	return nil
}
