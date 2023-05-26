package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/proveniencenft/kmsclitool/common"
	"github.com/proveniencenft/primesecrets/gf256"
)

var pass = []byte("aaaaaa")

func TestSplit(t *testing.T) {
	privkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("D:", privkey.D)
	s, err := gf256.SplitBytes(privkey.D.Bytes(), shares, threshold)

	kdf := "scrypt"
	encalg := "aes-128-ctr"
	keyfile := common.Keyfile{}
	crypto := &keyfile.Crypto
	crypto.Kdf = kdf
	crypto.Cipher = encalg
	err = common.EncryptAES(&keyfile, s[0].Value, pass)
	if err != nil {
		fmt.Println(err)
	}
	shareRaw, err := common.DecryptAES(&keyfile, pass)

	privDBytes, err := gf256.RecoverBytes(s[2:4])
	if err != nil {
		fmt.Println(err)
	}
	privKD := new(big.Int).SetBytes(privDBytes)
	fmt.Println("D:", privKD)

}
