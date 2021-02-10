package steganography

import "image"

const headerSize = 32

// CanFit checks if it is possible to put the goal image inside the target file
func CanFit(im image.Image, contents []byte, lsb int) bool {
	bitsToHide := len(contents) * 8

	maxes := im.Bounds().Max
	pixels := maxes.X * maxes.Y
	availiableBits := (pixels * 3 * lsb) - headerSize

	return availiableBits >= bitsToHide
}
