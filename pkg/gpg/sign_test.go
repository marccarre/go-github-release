package gpg_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/marccarre/go-github-release/pkg/gpg"

	"github.com/stretchr/testify/assert"
)

func TestArmoredDetachSign(t *testing.T) {
	signer, err := gpg.NewSigner("test-key.asc", "s3cr3t")
	assert.NoError(t, err)

	tmpFile, err := ioutil.TempFile("", "TestArmoredDetachSign")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	signature, err := signer.ArmoredDetachSign(tmpFile.Name())
	assert.NoError(t, err)
	signatureInfo, err := os.Stat(signature)
	assert.NoError(t, err)
	assert.False(t, signatureInfo.IsDir())
	assert.Equal(t, filepath.Base(tmpFile.Name())+".asc", signatureInfo.Name())
}
