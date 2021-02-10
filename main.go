package main

import (
	"fmt"
	"os"
	"steganographer/args"
)

func main() {
	imageFile, hiddenFile, lsb := args.GetFlags()

	im, toHide, err := args.ValidateArguments(imageFile, hiddenFile, lsb)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(im, toHide)
}
