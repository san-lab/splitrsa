package main

import (
	"fmt"
	"math/big"

	"github.com/manifoldco/promptui"
	"github.com/proveniencenft/primesecrets/gf256"
)

// TODO PedersenVSS
func VerifyShares() error {
	shares := make([]gf256.Share, 0)
	for {
		prompt := promptui.Select{
			Label: "SSS",
			Items: []string{"Input another share", "Verify", "EXIT"},
		}
		_, it, _ := prompt.Run()
		switch it {
		case "Input another share":

			// Pick files and prompt for passwords

			shareWrapper, err := ReadShare()
			if err != nil {
				fmt.Println(err)
			}

			pass := []byte(PromptForPassword("Password"))
			keyfile := shareWrapper.Keyfile
			err = DecryptKeyfile(&keyfile, pass)
			if err != nil {
				fmt.Println(err)
			}
			shareRaw := keyfile.Plaintext
			if err != nil {
				fmt.Println("Wrong password")
				return err // Handle better than this!
			}

			share := gf256.Share{Point: shareWrapper.Idx, Value: shareRaw, Degree: byte(shareWrapper.T - 1)}
			//fmt.Println(share)
			shares = append(shares, share)

		case "Verify":
			if len(shares) == 0 {
				fmt.Println("not enough shares!")
				break
			}
			privDBytes, err := gf256.RecoverBytes(shares)
			if err != nil {
				fmt.Println(err)
				break
			} else {
				new(big.Int).SetBytes(privDBytes)
				fmt.Println("All good!")
				//fmt.Println(D)
				// TODO calc the rest of the parameters to verify against pubKey
				break
			}

		case "EXIT":
			return nil
		}

	}
}
