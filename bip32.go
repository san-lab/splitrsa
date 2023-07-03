package main

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/proveniencenft/kmsclitool/common"
	"github.com/tyler-smith/go-bip32"
)

var currentKey *bip32.Key

const setb32 = "Set key"
const privhex = "Key in hex"
const chainhex = "Chaincode in hex"
const derive32 = "Derive"

func BIP32() {
	for {
		var keystring string
		var items []string
		if currentKey == nil {
			keystring = ">>KEY NOT SET<<"
			items = []string{setb32, up}
		} else {
			keystring = currentKey.B58Serialize()
			items = []string{setb32, privhex, chainhex, derive32, up}
		}
		sel := promptui.Select{Label: fmt.Sprintf("BIP32 for %s", keystring)}

		sel.Items = items
		_, fn, _ := sel.Run()
		switch fn {
		case up:
			return
		case setb32:
			SetB32Key()
		case privhex:
			fmt.Println(hex.EncodeToString(currentKey.Key))
		case chainhex:
			fmt.Println(hex.EncodeToString(currentKey.ChainCode))
		case derive32:
			DeriveBIP32()
		}

	}

}

func SetB32Key() {
	pr := promptui.Prompt{Label: "Input base58 encoded key"}
	b58, err := pr.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	ReadB58(b58)
}

func ReadB58(b58 string) error {
	key, err := bip32.B58Deserialize(b58)
	if err != nil {
		return err
	}
	currentKey = key
	return nil
}

func DeriveBIP32() {
	if currentKey == nil {
		fmt.Println("Key not set")
		return
	}

	path := GetDerivationPath()
	prettypath := "/"
	nextKey := currentKey
	var err error
	for _, idx := range path {
		nextKey, err = nextKey.NewChildKey(idx)
		if err != nil {
			fmt.Println(err)
			return
		}
		prettypath += fmt.Sprintf("%v/", idx)

	}
	//TODO handle nil
	fmt.Printf("Derivation path:\t  %s\n", prettypath)
	fmt.Printf("Derived key in base58:\t  %s\n", nextKey.B58Serialize())
	fmt.Printf("Derived key in hex:\t  %s\n", hex.EncodeToString(nextKey.Key))
	fmt.Printf("Ethereum address:\t  %s\n", toAddress(nextKey.Key))

}

var pathregex = regexp.MustCompile("[0-9]+")

func GetDerivationPath() []uint32 {
	path := []uint32{}
	pr := promptui.Prompt{Label: "Derivation path (/x/y/z/..)"}
	rawpath, _ := pr.Run()
	vs := pathregex.FindAllString(rawpath, -1)
	for _, ns := range vs {
		u64, err := strconv.ParseUint(ns, 10, 0)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		path = append(path, uint32(u64))
	}

	return path
}

func toAddress(privkey []byte) string {
	return common.CRCAddressFromPub(common.Scalar2Pub(privkey))

}
