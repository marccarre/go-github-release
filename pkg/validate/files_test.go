package validate_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/marccarre/go-github-release/pkg/validate"
	"github.com/stretchr/testify/assert"
)

func TestFilesNonExistingFile(t *testing.T) {
	assert.EqualError(t, validate.Files([]string{"/path/that/does/not/exist"}), "stat /path/that/does/not/exist: no such file or directory")
}

func TestFilesExistingFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "TestFilesExistingFile")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	assert.NoError(t, validate.Files([]string{tmpfile.Name()}))
}

func TestFileNonExistingFile(t *testing.T) {
	assert.EqualError(t, validate.File("/path/that/does/not/exist"), "stat /path/that/does/not/exist: no such file or directory")
}

func TestFileExistingFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "TestFileExistingFile")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	assert.NoError(t, validate.File(tmpfile.Name()))
}
