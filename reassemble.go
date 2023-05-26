package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/proveniencenft/kmsclitool/common"
	"github.com/proveniencenft/primesecrets/gf256"
)

const another = "Another share"
const assemble = "Assemble!"

func ReassemblePrivateKey() error {
	shares := make([]gf256.Share, 0)
	// Pick files and prompt for passwords
	for {
		shareWrapper, err := ReadShare()
		if err != nil {
			fmt.Println(err)
		}
		//b, _ := json.MarshalIndent(shareWrapper, " ", " ")
		//os.WriteFile("test-file", b, 0644)

		pass := []byte(PromptForPassword("Password"))
		keyfile := shareWrapper.Keyfile
		fmt.Println(keyfile.Crypto.KdfScryptParams.N)
		err = DecryptKeyfile(&keyfile, pass)
		if err != nil {
			fmt.Println(err)
		}
		shareRaw := keyfile.Plaintext
		//shareRaw, err := common.DecryptAES(&(shareWrapper.Keyfile), pass)
		if err != nil {
			fmt.Println("Wrong password")
			return err // Handle better than this!
		}

		share := gf256.Share{Point: shareWrapper.Idx, Value: shareRaw, Degree: byte(shareWrapper.T - 1)}
		//fmt.Println(share)
		shares = append(shares, share)
		if len(shares) == shareWrapper.T {
			privDBytes, err := gf256.RecoverBytes(shares)
			if err != nil {
				fmt.Println(err)
				break
			}
			D := new(big.Int).SetBytes(privDBytes)
			fmt.Println(D)
			break
		}
	}
	return nil
}

func ReadShare() (*Wrapper, error) {
	filename := "not-existing-file" // do better...
	for !fileExists(filename) {
		filename = PromptForString(fmt.Sprintf("File name of share"), "")
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	file := &Wrapper{}
	err = json.Unmarshal(data, file)

	file.Keyfile.UnmarshalKdfJSON()

	return file, err
}

func fileExists(name string) bool {
	path, err := os.Getwd()
	file := path + "/" + name
	_, err = os.Stat(file)
	fmt.Println(file)
	return err == nil
}

func DecryptKeyfile(kf *common.Keyfile, pass []byte) (err error) {
	key, err := kf.KeyFromPass(pass)
	if err != nil {
		return
	}
	fmt.Println("Verifying MAC...")
	err = kf.VerifyMAC(key)
	if err != nil {
		return
	}
	kf.Plaintext, err = common.Decrypt(kf, key)
	return
}
