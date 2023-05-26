package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/proveniencenft/kmsclitool/common"
	"github.com/proveniencenft/primesecrets/gf256"
)

var pattern = "ShamShare"

var threshold = 2
var shares = 4

func GenerateShares() error {
	klen := big.NewInt(1)
	for !isValidLength(klen) {
		klen = PromptForNumber("RSA key length?", "4096")
	}

	privkey, err := rsa.GenerateKey(rand.Reader, int(klen.Int64()))
	if err != nil {
		return err
	}
	//degree := threshold - 1
	filepat := "Splitkey"
	fmt.Println("D:", privkey.D)
	s, err := gf256.SplitBytes(privkey.D.Bytes(), shares, threshold)

	xuuid, err := uuid.NewUUID()
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

	return nil
}

type Wrapper struct {
	Keyfile common.Keyfile
	N       *big.Int
	E       int
	T       int
	Idx     byte
	ID      string
	L       *big.Int
}

func PromptForNumber(label, def string) *big.Int {
	pr := promptui.Prompt{Label: label, Default: def}
	v := new(big.Int)
	for {
		res, _ := pr.Run()
		_, ok := v.SetString(res, 10)
		if ok {
			return v
		}
		pr.Default = res
	}

}

func PromptForString(label, def string) string {
	pr := promptui.Prompt{Label: label, Default: def}
	res, err := pr.Run()
	if err != nil {
		fmt.Println(err)
		return def
	}
	return res
}

func PromptForPassword(label string) string {
	pr := promptui.Prompt{Label: label, Mask: '*'}
	res, err := pr.Run()
	if err != nil {
		fmt.Println(err)
		return "pass" // Default pass, dont show in prompt because its confusing with the mask
	}
	return res

}

func isValidLength(l *big.Int) bool {
	x := new(big.Int)
	x.Mod(l, big.NewInt(1024))
	return x.Cmp(big.NewInt(0)) == 0
}
