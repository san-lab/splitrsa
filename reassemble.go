package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/proveniencenft/kmsclitool/common"
	"github.com/proveniencenft/primesecrets/gf256"
)

const another = "Another share"
const assemble = "Assemble!"

func ReassemblePrivateKey() error {
	shares := make([][]byte, 0)
	// Pick files and prompt for passwords
	for {
		prompt := promptui.Select{
			Label: "SSS",
			Items: []string{another, assemble},
		}
		_, it, _ := prompt.Run()
		switch it {
		case another:
			i := 0
			filename := ""
			for !fileExists(filename) {
				filename = PromptForString(fmt.Sprintf("File name of share n%v", i), "")
			}
			pass := []byte(PromptForPassword("File password"))
			kdf := "scrypt"
			encalg := "aes-256-ctr"
			keyfile := common.Keyfile{}
			crypto := &keyfile.Crypto
			crypto.Kdf = kdf
			crypto.Cipher = encalg
			share, err := common.DecryptAES(&keyfile, pass)
			if err != nil {
				fmt.Println("Wrong password")
				return err // Handle better than this!
			}
			shares = append(shares, share)

		case "Assemble!":
			// Assemble
			gf256.RecoverBytes(shares)
		}

	}

}

func fileExists(name string) bool {
	path, err := os.Getwd()
	file := path + name
	_, err = os.Stat(file)
	return err == nil
}
