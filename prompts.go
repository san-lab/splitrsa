package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const up = "Up"
const gen = "Generate split"
const verify = "Verify shares"
const reass = "Reassemble key"

func TopUI() {
	for {
		prompt := promptui.Select{
			Label: "SSS",
			Items: []string{gen, verify, reass, "EXIT"},
		}
		_, it, _ := prompt.Run()
		switch it {
		case gen:
			GenerateShares()
		case verify:
			NoI()
		case reass:
			ReassemblePrivateKey()
		case "EXIT":
			return
		}

	}
}

func NoI() {
	fmt.Println("Not implemented")
}
