package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func a() {
	fmt.Println("b")
}

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
			Type:  "PUBLIC KEY",
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
			Type:  "PUBLIC KEY",
			Bytes: b,
		},
	)
	return pemdata, nil
}

func SavePub(key *rsa.PublicKey, filename string) {
	os.WriteFile(filename, PubPem(key), 0644)
	return
}

func SavePubPKIX(key *rsa.PublicKey, filename string) {
	b, err := PubPemPKIX(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	os.WriteFile(filename, b, 0644)
	return
}

func SavePriv(key *rsa.PrivateKey, filename string) {
	e := os.WriteFile(filename, PrivPem(key), 0644)
	if e != nil {
		fmt.Println(e)
	}
	return
}

func ReadPubPEM() *rsa.PublicKey {
	fileList, err := FindFilesWithExtension(".pem")
	prompt := promptui.Select{
		Label: "SSS",
		Items: fileList,
	}
	_, it, _ := prompt.Run()

	pubPEMData, err := os.ReadFile(it)
	if err != nil {
		return nil
	}
	fmt.Println(string(pubPEMData))
	block, _ := pem.Decode(pubPEMData)
	if block == nil || block.Type != "PUBLIC KEY" {
		fmt.Println("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pub)
	return pub.(*rsa.PublicKey)
}
