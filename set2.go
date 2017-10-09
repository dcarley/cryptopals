package cryptopals

import (
	"bytes"
	"math"
)

// PKCS7Padding adds PKCS#7 padding to a byte slice.
// https://en.wikipedia.org/wiki/Padding_(cryptography)#PKCS7
func PKCS7Padding(text []byte, blockSize int) []byte {
	factor := math.Ceil(float64(len(text)) / float64(blockSize))
	newSize := int(factor) * blockSize
	padSize := newSize - len(text)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)

	return append(text, pad...)
}
