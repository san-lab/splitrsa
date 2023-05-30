package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
)

func TestListFiles(t *testing.T) {
	fileList := make([]string, 0)
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fmt.Println(file.Name())
		fmt.Println(file.Name()[len(file.Name())-3:])
		if file.Name()[len(file.Name())-3:] == ".go" {
			fileList = append(fileList, file.Name())
		}
	}
	fmt.Println(fileList)
}

func TestPubPrivMatch(t *testing.T) {
	privPEMData, err := os.ReadFile("assembledPriv.pem")
	if err != nil {
		fmt.Println(err)
	}
	block, _ := pem.Decode(privPEMData)

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(priv)
	pubPem := string(PubPem(&priv.PublicKey))
	fmt.Println(pubPem)
	fmt.Println("--------------")

	privPEMData, err = os.ReadFile(".privkey.1685443285.pem")
	if err != nil {
		fmt.Println(err)
	}
	block, _ = pem.Decode(privPEMData)

	priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(priv)
	pubPem = string(PubPem(&priv.PublicKey))
	fmt.Println(pubPem)
}
