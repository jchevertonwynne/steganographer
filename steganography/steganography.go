package steganography

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

const headerSize = 32

// CanFit checks if it is possible to put the goal image inside the target file
func canFit(im image.Image, contents []byte, lsb int) error {
	maxes := im.Bounds().Max
	availiableBytes := ((maxes.X * maxes.Y * 3 * lsb) - headerSize) / 8

	if availiableBytes < len(contents) {
		return fmt.Errorf("only a maximum of %d bytes can be stored using %d least sig bits, target file has %d", availiableBytes/8, lsb, len(contents))
	}

	return nil
}

// Encode takes an image, some contents to hide and the number of least sig bits to use. Returns the bytes of the steganopraphised image
func Encode(im image.Image, contents []byte, lsb int) ([]byte, error) {
	err := canFit(im, contents, lsb)
	if err != nil {
		return nil, fmt.Errorf("failed to encode file: %v", err)
	}

	output := image.NewNRGBA(im.Bounds())
	draw.Draw(output, im.Bounds(), im, output.Bounds().Min, draw.Src)

	bits := getBits(contents)
	maxX := im.Bounds().Max.X
	shift := 0
	rgb := 0
	xInd := 0
	yInd := 0

	for _, currentBit := range bits {
		col := output.NRGBAAt(xInd, yInd)

		r := col.R
		g := col.G
		b := col.B
		a := col.A

		switch currentBit {
		case 0:
			switch rgb {
			case 0:
				r &= ^(1 << shift)
			case 1:
				g &= ^(1 << shift)
			case 2:
				b &= ^(1 << shift)
			default:
				panic("rgb out of index")
			}
		case 1:
			switch rgb {
			case 0:
				r |= 1 << shift
			case 1:
				g |= 1 << shift
			case 2:
				b |= 1 << shift
			default:
				panic("rgb out of index")
			}
		default:
			panic("bad bit value")
		}

		output.SetNRGBA(xInd, yInd, color.NRGBA{r, g, b, a})

		shift++
		if shift == lsb {
			shift = 0
			rgb++
			if rgb == 3 {
				rgb = 0
				xInd++
				if xInd == maxX {
					xInd = 0
					yInd++
				}
			}
		}
	}

	res := new(bytes.Buffer)
	err = png.Encode(res, output)
	if err != nil {
		return nil, fmt.Errorf("failed to encode file: %v", err)
	}

	return res.Bytes(), nil
}

func getBits(contents []byte) []int {
	bits := make([]int, 0, headerSize)

	length := len(contents)
	for i := 0; i < headerSize; i++ {
		bits = append(bits, length&1)
		length >>= 1
	}

	for _, b := range contents {
		for i := 0; i < 8; i++ {
			bits = append(bits, int(b&1))
			b >>= 1
		}
	}

	return bits
}

// Decode attempts to retrive data from a steganographised image
func Decode(im image.Image, lsb int) ([]byte, error) {
	imNGGBA := image.NewNRGBA(im.Bounds())
	draw.Draw(imNGGBA, im.Bounds(), im, imNGGBA.Bounds().Min, draw.Src)

	contentLength, err := getContentsSize(im, lsb)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %v", err)
	}

	var result []byte

	maxX := im.Bounds().Max.X
	shift := 0
	rgb := 0
	xInd := 0
	yInd := 0

	shift = headerSize
	rgb = shift / lsb
	shift %= lsb
	xInd = rgb / 3
	rgb %= 3
	yInd = xInd / maxX
	xInd %= maxX

	for i := 0; i < contentLength; i++ {
		var currentBit int
		for j := 0; j < 8; j++ {
			col := imNGGBA.NRGBAAt(xInd, yInd)
			r := col.R
			g := col.G
			b := col.B

			switch rgb {
			case 0:
				if r&(1<<shift) != 0 {
					currentBit |= 1 << j
				}
			case 1:
				if g&(1<<shift) != 0 {
					currentBit |= 1 << j
				}
			case 2:
				if b&(1<<shift) != 0 {
					currentBit |= 1 << j
				}
			default:
				panic("rgb out of index")
			}

			shift++
			if shift == lsb {
				shift = 0
				rgb++
				if rgb == 3 {
					rgb = 0
					xInd++
					if xInd == maxX {
						xInd = 0
						yInd++
					}
				}
			}
		}

		result = append(result, byte(currentBit))
	}

	return result, nil
}

func getContentsSize(im image.Image, lsb int) (int, error) {
	bounds := im.Bounds().Max

	totalStorage := bounds.X * bounds.Y * 3 * lsb

	if totalStorage < 32 {
		return 0, fmt.Errorf("image is too small to get content size: only %d bits availiable", totalStorage)
	}

	maxX := bounds.X
	shift := 0
	rgb := 0
	xInd := 0
	yInd := 0

	result := 0

	for i := 0; i < headerSize; i++ {
		r, g, b, _ := im.At(xInd, yInd).RGBA()

		switch rgb {
		case 0:
			if r&(1<<shift) != 0 {
				result |= (1 << i)
			}
		case 1:
			if g&(1<<shift) != 0 {
				result |= (1 << i)
			}
		case 2:
			if b&(1<<shift) != 0 {
				result |= (1 << i)
			}
		default:
			panic("rgb out of index")
		}

		shift++
		if shift == lsb {
			shift = 0
			rgb++
			if rgb == 3 {
				rgb = 0
				xInd++
				if xInd == maxX {
					xInd = 0
					yInd++
				}
			}
		}
	}

	if result*8+32 > totalStorage {
		return 0, fmt.Errorf("image has invalid size: %d bytes is too many to fit in using %d least sig bits", result, lsb)
	}

	return result, nil
}
