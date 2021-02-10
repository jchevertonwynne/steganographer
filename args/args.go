package args

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"steganographer/steganography"
)

// GetFlags gets the image file name, the file to hide and the number of least sig bits to use in hiding
func GetFlags() (string, string, int) {
	imageFile := flag.String("i", "", "image file to hide data in")
	hiddenFile := flag.String("t", "", "file to hide within image")
	lsb := flag.Int("b", 1, "number of least significant bits to use for steganography")
	flag.Parse()

	return *imageFile, *hiddenFile, *lsb
}

// ValidateArguments ensures that the files exist and are of the correct type and size
func ValidateArguments(imageFilename, hiddenFilename string, lsb int) (image.Image, []byte, error) {
	if imageFilename == "" {
		return nil, nil, errors.New("please enter an image file name")
	}

	if hiddenFilename == "" {
		return nil, nil, errors.New("please enter a target file name")
	}

	imageFile, err := os.Open(imageFilename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open %q: %v", imageFilename, err)
	}
	defer imageFile.Close()

	loadedPNG, err := png.Decode(imageFile)
	if err != nil {
		return nil, nil, errors.New("image must be a PNG")
	}

	hiddenFile, err := os.Open(hiddenFilename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open %q: %v", hiddenFilename, err)
	}
	defer hiddenFile.Close()

	hiddenContents, err := ioutil.ReadAll(hiddenFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read hidden file: %v", err)
	}

	if lsb < 1 || lsb > 8 {
		return nil, nil, errors.New("least sig bits must be between 1 and 8")
	}

	if !steganography.CanFit(loadedPNG, hiddenContents, lsb) {
		return nil, nil, fmt.Errorf("image %q is not large enough to fit %q with %d least sig bit changes", imageFilename, hiddenFilename, lsb)
	}

	return loadedPNG, hiddenContents, nil
}
