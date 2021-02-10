package main

import (
	"fmt"
	"os"
	"steganographer/args"
	"steganographer/steganography"
)

func main() {
	imageFile, hiddenFile, lsb, decode := args.GetFlags()

	im, toHide, err := args.CheckValid(imageFile, hiddenFile, lsb, decode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if decode {
		output := steganography.Decode(im, lsb)
		fmt.Println(string(output))
	} else {
		output, err := steganography.Encode(im, toHide, lsb)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(string(output))
	}
}
