package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func PrivPem(privk *rsa.PrivateKey) []byte {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privk),
		},
	)
	return pemdata
}

func PubPem(pubk *rsa.PublicKey) []byte {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubk),
		},
	)
	return pemdata
}

func PubPemPKIX(pubk *rsa.PublicKey) ([]byte, error) {
	b, e := x509.MarshalPKIXPublicKey(pubk)
	if e != nil {
		return nil, e
	}
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: b,
		},
	)
	return pemdata, nil
}
