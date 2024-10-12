package extractembed

import (
	"embed"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test1/*.txt test1/test2/a.txt
var embedFiles embed.FS

func TestExtract(t *testing.T) {
	err := Extract("tmp/test-extract", &embedFiles)
	if err != nil {
		t.Fatal(err)
	}

	output, err := exec.Command("find", "tmp/test-extract").Output()
	if err != nil {
		t.Fatal(err)
	}

	excepted := `tmp/test-extract
tmp/test-extract/test1
tmp/test-extract/test1/b.txt
tmp/test-extract/test1/test2
tmp/test-extract/test1/test2/a.txt
tmp/test-extract/test1/a.txt
`

	assert.Equal(t, excepted, string(output))

	_, err = exec.Command("rm", "-rf", "tmp/test-extract").Output()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)

}
