package gpg

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

// Signer encapsulates low-level openpgp
type Signer struct {
	signer *openpgp.Entity
}

// NewSigner creates a new signer based on the the provided private key.
// If the private key is encrypted, the provided (optional) password will be
// used to decrypt the private key.
func NewSigner(privateKeyPath, password string) (*Signer, error) {
	key, err := os.Open(filepath.Clean(privateKeyPath))
	if err != nil {
		return nil, fmt.Errorf("failed to open GPG private key: %s", err)
	}
	defer key.Close()
	signer, err := openpgp.ReadEntity(packet.NewReader(key))
	if err != nil {
		return nil, fmt.Errorf("failed to read GPG private key: %s", err)
	}
	if signer.PrivateKey.Encrypted {
		if err = signer.PrivateKey.Decrypt([]byte(password)); err != nil {
			return nil, fmt.Errorf("failed to decrypt GPG private key: %s", err)
		}
	}
	return &Signer{
		signer: signer,
	}, nil
}

// ArmoredDetachSign generates an armored signature of the provided file, as a
// separate file, and returns the path to the signature file.
func (s Signer) ArmoredDetachSign(filePath string) (string, error) {
	r, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return "", err
	}
	defer r.Close()
	signaturePath := filePath + ".asc"
	w, err := os.OpenFile(signaturePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return "", err
	}
	defer w.Close()
	if err = openpgp.ArmoredDetachSign(w, s.signer, r, &packet.Config{}); err != nil {
		return "", err
	}
	return signaturePath, nil
}
