package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	RSAPrivateKeyType = "RSA PRIVATE KEY"
	RSAPublicKeyType  = "RSA PUBLIC KEY"

	PublicFileName  = "private_key.pem"
	PrivateFileName = "public_key.pem"

	BitSizeKey = 4096
	FilePerm   = 0600
)

func main() {
	// создаём новый RSA-ключ длиной 4096 бит
	var (
		outPrivateKeyOutFilePath string
		outPublicKeyOutFilePath  string
	)
	flag.StringVar(&outPrivateKeyOutFilePath, "private", ".", "private file pathout")
	flag.StringVar(&outPublicKeyOutFilePath, "public", ".", "public file path out")
	flag.Parse()

	privateKey, err := rsa.GenerateKey(rand.Reader, BitSizeKey)
	pubKey := &privateKey.PublicKey
	if err != nil {
		log.Fatal(err)
	}

	var privateKeyPEM bytes.Buffer
	if err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  RSAPrivateKeyType,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		log.Fatal(err)
	}

	var publickKeyPEM bytes.Buffer
	if err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  RSAPublicKeyType,
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	}); err != nil {
		log.Fatal(err)
	}

	if err = checkKeys(privateKey, pubKey); err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(outPrivateKeyOutFilePath, PrivateFileName), privateKeyPEM.Bytes(), FilePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(outPublicKeyOutFilePath, PublicFileName), publickKeyPEM.Bytes(), FilePerm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Сертификат и приватный ключ успешно сохранены в файлы.")
}

func checkKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) error {
	message := []byte("Это секретное сообщение.")

	encryptedMessage, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Зашифрованное сообщение: %x\n", encryptedMessage)

	// Пример расшифровки
	decryptedMessage, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedMessage, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Расшифрованное сообщение: %s\n", decryptedMessage)
	if !bytes.Equal(decryptedMessage, message) {
		return errors.New("failed to check keys")
	}
	return nil
}
