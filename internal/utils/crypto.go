package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

type PublicKey struct {
	pub any
}

type PrivateKey struct {
	priv any
}

// LoadPublicKey loads a public key from the specified file path.
//
// It takes a filePath string as a parameter and returns a *PublicKey and an error.
func LoadPublicKey(filePath string) (*PublicKey, error) {
	if filePath == "" {
		return nil, nil
	}
	publicKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &PublicKey{pub: publicKey}, nil
}

// LoadPrivateKey loads a private key from the specified file path.
//
// It takes a filePath string as a parameter and returns a *PrivateKey and an error.
func LoadPrivateKey(filePath string) (*PrivateKey, error) {
	if filePath == "" {
		return nil, nil
	}
	privateKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{priv: privateKey}, nil
}

// Encrypt encrypts the given data using the public key.
//
// It takes a string parameter 'data' which represents the data to be encrypted.
// It returns a string which represents the encrypted data and an error if any.
func (pub *PublicKey) Encrypt(data []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub.pub.(*rsa.PublicKey), data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// Decrypt decrypts the given data using the private key.
//
// The data parameter is the ciphertext to be decrypted.
// It returns the plaintext string and an error if decryption fails.
func (priv *PrivateKey) Decrypt(data []byte) ([]byte, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv.priv.(*rsa.PrivateKey), data)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
