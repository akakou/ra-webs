package extractembed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListEmbedFiles(t *testing.T) {
	fil, err := listEmbedFiles(".", &embedFiles)
	assert.NoError(t, err)

	assert.Equal(t, []string{"test1/a.txt", "test1/b.txt", "test1/test2/a.txt"}, fil)
}

func TestReadEmbedFile(t *testing.T) {
	res, err := readEmbedFile("test1/a.txt", &embedFiles)
	assert.NoError(t, err)

	assert.Equal(t, "Hello :)", string(res))
}
