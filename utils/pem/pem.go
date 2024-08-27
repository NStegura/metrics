package pem

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	RSAPrivateKeyType = "RSA PRIVATE KEY"
	RSAPublicKeyType  = "RSA PUBLIC KEY"
)

func ReadPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file, %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != RSAPrivateKeyType {
		return nil, errPrivateKeyType
	}

	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key, %w", err)
	}
	return pk, nil
}

func ReadPublicKey(filePath string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file, %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != RSAPublicKeyType {
		return nil, errPublicKeyType
	}

	pubInterface, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse publick key")
	}

	return pubInterface, nil
}
