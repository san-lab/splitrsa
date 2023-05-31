package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/proveniencenft/kmsclitool/common"
	"github.com/proveniencenft/primesecrets/gf256"
)

func ReassemblePrivateKey() error {
	shares := make([]gf256.Share, 0)
	// Pick files and prompt for passwords
	for {
		shareWrapper, err := ReadShare()
		if err != nil {
			if err.Error() == "exited" {
				break // Here user just wants to go back
			}
			fmt.Println(err)
			continue // Otherwise, its better not to reset the shares already provided
		}

		pass := []byte(PromptForPassword("Password"))
		keyfile := shareWrapper.Keyfile
		//fmt.Println(keyfile.Crypto.KdfScryptParams.N)
		err = DecryptKeyfile(&keyfile, pass)
		if err != nil {
			fmt.Println(err)
			continue
		}
		shareRaw := keyfile.Plaintext
		if err != nil {
			fmt.Println("Wrong password")
			continue
		}

		share := gf256.Share{Point: shareWrapper.Idx, Value: shareRaw, Degree: byte(shareWrapper.T - 1)}
		shares = append(shares, share)

		if len(shares) >= shareWrapper.T {
			privDBytes, err := gf256.RecoverBytes(shares) // TODO FIX BUG? WHEN REPEATED SHARES
			fmt.Println("err:", err)
			if err != nil {
				if err.Error() == "Not enough shares" {
					// This means user has input the same shares more than once,
					// which is fine but dont want to reset the shares already provided
					continue
				}
				break // Idk about the rest of the errors...
			}
			D := new(big.Int).SetBytes(privDBytes)
			pubK, err := ReadPubKey()
			if err != nil {
				return err
			}
			prk1, err := D2PrivKey(D, pubK)
			err = prk1.Validate() // This checks if the pubK contained is consistent with the private values
			if err != nil {
				fmt.Println(err)
				return err
			}
			SavePriv(prk1, "assembledPriv.pem")
			fmt.Println("Key reassembled and saved!")
			break
		}
	}
	return nil
}

func ReadShare() (*Wrapper, error) {
	fileList, err := FindFilesWithExtension(".json")
	it := PromptFromList(fileList, "Pick a share file")
	if it == "EXIT" {
		return nil, errors.New("exited")
	}

	data, err := os.ReadFile(it)
	if err != nil {
		return nil, err
	}
	file := &Wrapper{}
	err = json.Unmarshal(data, file)

	file.Keyfile.UnmarshalKdfJSON() // Set share metadata

	return file, err
}

func ReadPubKey() (*rsa.PublicKey, error) {
	fileList, err := FindFilesWithExtension(".pem")
	pubKPemFile := PromptFromList(fileList, "Pick a public key file")
	pubKPemBytes, err := os.ReadFile(pubKPemFile)
	if err != nil {
		return nil, err
	}
	pubK, err := ParsePubPem(pubKPemBytes)
	return pubK, err
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
