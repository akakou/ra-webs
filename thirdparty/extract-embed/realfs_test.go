package extractembed

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveFile(t *testing.T) {
	err := saveFile("tmp.txt", []byte("Hello"), "./tmp/save-file")
	assert.NoError(t, err)

	a, err := os.ReadFile("./tmp/save-file/tmp.txt")
	assert.NoError(t, err)

	assert.Equal(t, "Hello", string(a))

	err = os.Remove("./tmp/save-file/tmp.txt")
	assert.NoError(t, err)

	err = os.Remove("./tmp/save-file")
	assert.NoError(t, err)
}
