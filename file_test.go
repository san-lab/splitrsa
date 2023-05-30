package main

import (
	"fmt"
	"os"
	"testing"
)

func TestListFiles(t *testing.T) {
	fileList := make([]string, 0)
	files, _ := os.ReadDir(".")
	for _, file := range files {
		fmt.Println(file.Name())
		fmt.Println(file.Name()[len(file.Name())-3:])
		if file.Name()[len(file.Name())-3:] == ".go" {
			fileList = append(fileList, file.Name())
		}
	}
	fmt.Println(fileList)
}
