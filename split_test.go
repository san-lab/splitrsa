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
	s, err := gf256.SplitBytes(privkey.D.Bytes(), 5, 3)
	if err != nil {
		t.Error(err)
	}

	b, err := gf256.RecoverBytes(s[2:])
	if err != nil {
		t.Error(err)
	}
	D1 := new(big.Int).SetBytes(b)
	if privkey.D.Cmp(D1) != 0 {
		t.Errorf("Wrong D recovered: /n%v/n%v", privkey.D, D1)
	}

}

func TestDecrypt(t *testing.T) {
	privkey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println(err)
	}
	pass := []byte("aaaaaa")
	kdf := "scrypt"
	encalg := "aes-256-ctr"
	keyfile := &common.Keyfile{}
	crypto := &keyfile.Crypto
	crypto.Kdf = kdf
	crypto.Cipher = encalg
	err = common.EncryptAES(keyfile, privkey.D.Bytes(), pass)
	if err != nil {
		t.Error(err)
	}
	key, err := keyfile.KeyFromPass(pass)
	if err != nil {
		t.Errorf("KDF fail: %s", err)
	}
	err = keyfile.VerifyMAC(key)
	if err != nil {
		t.Error("Wrong MAC")
	}
	keyfile.Decrypt(pass)
	b, err := common.Decrypt(keyfile, key)
	if err != nil {
		t.Error(err)
	}
	D1 := new(big.Int).SetBytes(b)
	if privkey.D.Cmp(D1) != 0 {
		t.Errorf("Wrong D recovered: \n%v\n%v", privkey.D, D1)
	}

}
