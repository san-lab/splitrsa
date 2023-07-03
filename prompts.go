package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
)

const rsaui = "RSA SSS"
const bip32ui = "BIP32"

var exit = CFormat("exit", &Yellow, nil)

func TopUI() {
	sel := promptui.Select{Label: Lbl("Select function:"), Items: []string{rsaui, bip32ui, exit}}
	for {
		_, fn, _ := sel.Run()
		switch fn {
		case rsaui:
			RSAUI()
		case bip32ui:
			BIP32()
		case exit:
			return
		}
	}
}

func Lbl(m string) string {
	return CFormat(m, &Cyan, &Bold)
}

var up = CFormat("up", &Yellow, &Underline)

const gen = "Generate split"
const verify = "Verify shares"
const reass = "Reassemble key"

func RSAUI() {
	for {
		prompt := promptui.Select{
			Label: Lbl("RSA SSS"),
			Items: []string{gen, verify, reass, up},
		}
		_, it, _ := prompt.Run()
		switch it {
		case gen:
			GenerateShares()
		case verify:
			VerifyShares()
		case reass:
			ReassemblePrivateKey()
		case up:
			return
		}

	}
}

func NoI() {
	fmt.Println("Not implemented")
}

func PromptForNumber(label string, def int) int {
	pr := promptui.Prompt{Label: label, Default: fmt.Sprint(def)}
	for {
		res, _ := pr.Run()
		v, err := strconv.Atoi(res)
		if err != nil {
			fmt.Println(err)
			return def
		}
		return v
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

func PromptFromList(list []string, label string) string {
	list = append(list, "EXIT")
	prompt := promptui.Select{
		Label: label,
		Items: list,
	}
	_, it, _ := prompt.Run()
	return it
}

func FindFilesWithExtension(ext string) ([]string, error) {
	fileList := make([]string, 0)
	files, err := os.ReadDir(".")
	for _, file := range files {
		if len(file.Name()) >= len(ext) && file.Name()[len(file.Name())-len(ext):] == ext {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList, err
}

func fileExists(name string) bool {
	path, err := os.Getwd()
	file := path + "/" + name
	_, err = os.Stat(file)
	fmt.Println(file)
	return err == nil
}
