package pem

import "errors"

var (
	errPrivateKeyType = errors.New("invalid private key format")
	errPublicKeyType  = errors.New("invalid public key format")
)
