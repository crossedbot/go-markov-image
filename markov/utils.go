package markov

import (
	"image/color"
)

const (
	DefaultCompression = 1
)

// encodeColor encodes a given color to a 32-bit unsigned integer.
func encodeColor(c color.Color) uint32 {
	c = compressColor(c, DefaultCompression)
	r, g, b, a := c.RGBA()
	return pack([4]uint8{uint8(r), uint8(g), uint8(b), uint8(a)})
}

// decodeColor decodes a given 32-bit unsigned integer to color.
func decodeColor(i uint32) color.Color {
	b := unpack(i)
	return color.RGBA{R: b[0], G: b[1], B: b[2], A: b[3]}
}

// Pack packs 4 8-bit unsigned integers into a single 32-bit unsigned integer.
func pack(b [4]uint8) uint32 {
	return uint32(b[0])<<24 |
		uint32(b[1])<<16 |
		uint32(b[2])<<8 |
		uint32(b[3])
}

// Unpack unpacks a single 32-bit unsigned integer int 4 8-bit unsigned
// integers.
func unpack(i uint32) [4]uint8 {
	return [4]uint8{
		uint8((i >> 24) & 255),
		uint8((i >> 16) & 255),
		uint8((i >> 8) & 255),
		uint8(i & 255),
	}
}

// colorToRGBA translates a color to its RGBA representation.
func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}

// CompressColor compresses a given color using the given threshold. RGBA values
// (R, G, B, or A) that are less than the threshold, will be reduced to 0, all
// other values are reduced to the closest multiple of the threshold; IE. given
// a R value of 22 and a threshold of 3, R will be compressed to 21.
func compressColor(c color.Color, threshold uint8) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: (uint8(r) / threshold) * threshold,
		G: (uint8(g) / threshold) * threshold,
		B: (uint8(b) / threshold) * threshold,
		A: (uint8(a) / threshold) * threshold,
	}
}
