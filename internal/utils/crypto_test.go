package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateKeys(publicKeyPath, privateKeyPath string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile(privateKeyPath, privateKeyPEM, 0644)
	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile(publicKeyPath, publicKeyPEM, 0644)
	if err != nil {
		panic(err)
	}
}

func TestCrypto(t *testing.T) {
	data := []byte("Hello, World!")
	publicKeyPath := "public.pem"
	privateKeyPath := "private.pem"
	generateKeys(publicKeyPath, privateKeyPath)

	publicKey, err := LoadPublicKey(publicKeyPath)
	assert.Nil(t, err)
	privateKey, err := LoadPrivateKey(privateKeyPath)
	assert.Nil(t, err)

	encryptedData, err := publicKey.Encrypt(data)
	assert.Nil(t, err)
	decryptedData, err := privateKey.Decrypt(encryptedData)
	assert.Nil(t, err)

	assert.Equal(t, data, decryptedData)
}
