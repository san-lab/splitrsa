package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/proveniencenft/kmsclitool/common"
	"github.com/proveniencenft/primesecrets/gf256"
)

var pattern = "ShamShare"

func GenerateShares() error {
	klen := 1
	for !isValidLength(klen) {
		klen = PromptForNumber("RSA key length?", 4096)
	}
	fmt.Printf("Generating %v-bit RSA key, please, wait...\n", klen)
	privkey, err := rsa.GenerateKey(rand.Reader, klen)
	if err != nil {
		return err
	}

	SavePriv(privkey, fmt.Sprintf(".privkey.%v.pem", time.Now().Unix()))

	fmt.Printf("RSA key generated\n")
	threshold := PromptForNumber("Quorum (threshold)?", 2)

	shares := PromptForNumber("Number of shares?", 4)
	//degree := threshold - 1
	filepat := "Splitkey"
	s, err := gf256.SplitBytes(privkey.D.Bytes(), shares, threshold)

	xuuid, err := uuid.NewUUID()
	fmt.Println("Share set ID:", xuuid.String())
	for _, sh := range s {
		filename := fmt.Sprintf("%s%v.json", filepat, sh.Point)
		filename = PromptForString(fmt.Sprintf("File name for share no %v", sh.Point), filename)
		pass := []byte(PromptForPassword("File password"))

		kdf := "scrypt"
		encalg := "aes-128-ctr"
		keyfile := common.Keyfile{}
		crypto := &keyfile.Crypto
		crypto.Kdf = kdf
		crypto.Cipher = encalg
		err := common.EncryptAES(&keyfile, sh.Value, pass)
		if err != nil {
			fmt.Println(err)
			return err
		}

		wrapper := Wrapper{}
		wrapper.N = privkey.PublicKey.N
		wrapper.E = privkey.PublicKey.E
		wrapper.Idx = sh.Point
		wrapper.T = threshold
		wrapper.Keyfile = keyfile
		wrapper.ID = xuuid.String()
		wrapper.L = klen

		b, _ := json.MarshalIndent(wrapper, " ", " ")
		os.WriteFile(filename, b, 0644)

	}
	pname := PromptForString("File name for the Public Key?", "pubkey")
	SavePubPKIX(&privkey.PublicKey, pname+"PKIX.pem")
	SavePub(&privkey.PublicKey, pname+".pem")

	fmt.Printf("Public key written to files:\n %s\n%s\n", pname+"PKIX.pem", pname+".pem")

	return nil
}

type Wrapper struct {
	Keyfile common.Keyfile
	N       *big.Int
	E       int
	T       int
	Idx     byte
	ID      string
	L       int
}

func isValidLength(l int) bool {
	return l%1024 == 0
}
