package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"math/big"
	"testing"
)

func TestGenerate(t *testing.T) {
	prk, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Error(err)
	}

	//SavePub(&prk.PublicKey, "testpub.yyy")
	//SavePriv(prk, "testpriv.yyy")

	prb := x509.MarshalPKCS1PrivateKey(prk)

	priv, err := x509.ParsePKCS1PrivateKey(prb)
	if err != nil {
		t.Error(err)
	}

	P, Q := Crack(priv.N, big.NewInt(int64(priv.E)), priv.D)

	prk1 := &rsa.PrivateKey{}
	prk1.N = prk.N
	prk1.D = prk.D
	prk1.PublicKey = prk.PublicKey
	prk1.Primes = []*big.Int{P, Q}
	prk1.Precompute()
	fmt.Println(prk.Equal(prk1))

	msg := []byte("Jasio Karuzela")
	ctx, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &prk.PublicKey, msg, nil)
	ptx, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, prk1, ctx, nil)
	if bytes.Compare(msg, ptx) != 0 {
		t.Error("RSA Decryption on recovered key failed")
	}
	fmt.Println("Recovered key decryption OK")

}
