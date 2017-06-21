package set1

import (
	"encoding/base64"
	"fmt"
)

// decToHex is used to lookup a single hex value in decimal.
var decToHex = map[byte]byte{
	0:  '0',
	1:  '1',
	2:  '2',
	3:  '3',
	4:  '4',
	5:  '5',
	6:  '6',
	7:  '7',
	8:  '8',
	9:  '9',
	10: 'a',
	11: 'b',
	12: 'c',
	13: 'd',
	14: 'e',
	15: 'f',
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

// HexToBase64 converts hexidecimal encoded text to base64.
func HexToBase64(text []byte) ([]byte, error) {
	raw, err := HexDecode(text)
	if err != nil {
		return []byte{}, err
	}

	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(b64, raw)

	return b64, nil
}
