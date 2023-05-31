package main

import (
	"fmt"
	"math/big"

	"github.com/proveniencenft/primesecrets/gf256"
)

// TODO PedersenVSS
func VerifyShares() error {
	shares := make([]gf256.Share, 0)
	for {
		it := PromptFromList([]string{"Input another share", "Verify"}, "Share verification")
		switch it {
		case "Input another share":
			shareWrapper, err := ReadShare("Reassembling an RSA key\nPlease, select the file with the first key share")
			if err != nil {
				fmt.Println(err)
				break
			}

			pass := []byte(PromptForPassword("Password"))
			keyfile := shareWrapper.Keyfile
			err = DecryptKeyfile(&keyfile, pass)
			if err != nil {
				fmt.Println(err)
				break
			}
			shareRaw := keyfile.Plaintext
			if err != nil {
				fmt.Println("Wrong password")
				break
			}

			share := gf256.Share{Point: shareWrapper.Idx, Value: shareRaw, Degree: byte(shareWrapper.T - 1)}
			shares = append(shares, share)

		case "Verify":
			if len(shares) == 0 {
				fmt.Println("Not enough shares!")
				break
			}
			privDBytes, err := gf256.RecoverBytes(shares)
			if err != nil {
				fmt.Println(err)
				break
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
			fmt.Println("All good!")
			break

		case "EXIT":
			return nil
		}

	}
}
