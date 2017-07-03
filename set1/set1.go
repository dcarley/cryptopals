package set1

import (
	"fmt"
)

// decToHex is used to lookup a single hex value in decimal.
var decToHex = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f',
}

// Encodes a decimal byte slice to a hexidecimal byte slice
func HexEncode(text []byte) []byte {
	// output is always twice the size of input
	out := make([]byte, len(text)*2)

	// based on this article:
	// https://learn.sparkfun.com/tutorials/hexadecimal#converting-tofrom-decimal
	for i, char := range text {
		secondDigit := decToHex[char%16]
		firstDigit := decToHex[char/16]
		out[i*2] = firstDigit
		out[i*2+1] = secondDigit
	}

	return out
}

// hexToDec is used to lookup a single decimal value in hex.
func hexToDec(char byte) (byte, error) {
	var val byte
	// subtract ASCII decimal codes so that:
	//   "0" becomes 0
	//   "a" and "A" become 10
	switch {
	case char >= '0' && char <= '9':
		val = char - '0'
	case char >= 'A' && char <= 'F':
		val = char - 'A' + 10
	case char >= 'a' && char <= 'f':
		val = char - 'a' + 10
	default:
		return 0, fmt.Errorf("invalid hex character: %s", []byte{char})
	}

	return val, nil
}

// HexDecode decodes a hexidecimal byte slice to a decimal byte slice
func HexDecode(text []byte) ([]byte, error) {
	if len(text)%2 != 0 {
		return []byte{}, fmt.Errorf("input must be an even size")
	}

	// output is always half the size of input
	out := make([]byte, len(text)/2)

	// based on this article:
	// https://learn.sparkfun.com/tutorials/hexadecimal#converting-tofrom-decimal
	for i := 0; i < len(out); i++ {
		firstDigit, err := hexToDec(text[i*2])
		if err != nil {
			return []byte{}, err
		}

		secondDigit, err := hexToDec(text[i*2+1])
		if err != nil {
			return []byte{}, err
		}

		out[i] = firstDigit*16 + secondDigit
	}

	return out, nil
}

// decToBase64 is used to lookup base64 characters using zero-indexed 6bit
// values
const decToBase64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// sixBitsMax is a 6bit value of all binary 1s:
//	- bin: 111111
//	- dec: 63
//	- hex: 0x3F
const sixBitsMax = 0x3F

// Base64Encode encodes a byte slice into a base64 byte slice
func Base64Encode(text []byte) []byte {
	// 4 bytes out for every 3 bytes in
	out := make([]byte, len(text)/3*4)

	for i := 0; i < len(text); i += 3 {
		// convert three 8bit characters into one 24bit value by:
		// - converting from byte (int8) to int32
		// - bit shifting to the left, to pad least significant bits, so that they don't overlap
		// - bit ORing to combine the values into one
		combined := int32(text[i]) << 16
		combined |= int32(text[i+1]) << 8
		combined |= int32(text[i+2]) << 0

		// split the 24bit value into four 6bit values by:
		// - bit shifting to the right, so the least significant 6bits are what we want
		// - bit ANDing against 6bits of binary 1s to extract the least significant 6bits
		// - looking up the appropriate base64 character for the value
		outIndex := i / 3 * 4
		out[outIndex] = decToBase64[combined>>18&sixBitsMax]
		out[outIndex+1] = decToBase64[combined>>12&sixBitsMax]
		out[outIndex+2] = decToBase64[combined>>6&sixBitsMax]
		out[outIndex+3] = decToBase64[combined>>0&sixBitsMax]
	}

	return out
}

// HexToBase64 converts hexidecimal encoded text to base64.
func HexToBase64(text []byte) ([]byte, error) {
	raw, err := HexDecode(text)
	if err != nil {
		return []byte{}, err
	}

	return Base64Encode(raw), nil
}
