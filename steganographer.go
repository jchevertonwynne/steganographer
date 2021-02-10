package main

import (
	"fmt"
	"os"
	"steganographer/args"
	"steganographer/steganography"
)

func main() {
	imageFile, hiddenFile, lsb, decode := args.GetFlags()

	im, toHide, err := args.Validate(imageFile, hiddenFile, lsb, decode)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if decode {
		output, err := steganography.Decode(im, lsb)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		fmt.Println(string(output))
	} else {
		output, err := steganography.Encode(im, toHide, lsb)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		fmt.Println(string(output))
	}
}
