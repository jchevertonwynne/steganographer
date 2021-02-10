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
func CanFit(im image.Image, contents []byte, lsb int) bool {
	bitsToHide := len(contents) * 8

	maxes := im.Bounds().Max
	pixels := maxes.X * maxes.Y
	availiableBits := (pixels * 3 * lsb) - headerSize

	return availiableBits >= bitsToHide
}

// Encode takes an image, some contents to hide and the number of least sig bits to use. Returns the bytes of the steganopraphised image
func Encode(im image.Image, contents []byte, lsb int) ([]byte, error) {
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
	err := png.Encode(res, output)
	if err != nil {
		return nil, fmt.Errorf("failed to hide contents in file: %v", err)
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
func Decode(im image.Image, lsb int) []byte {
	imNGGBA := image.NewNRGBA(im.Bounds())
	draw.Draw(imNGGBA, im.Bounds(), im, imNGGBA.Bounds().Min, draw.Src)

	var result []byte

	contentLength := getContentsSize(im, lsb)

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

	return result
}

func getContentsSize(im image.Image, lsb int) int {
	maxX := im.Bounds().Max.X
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

	return result
}
