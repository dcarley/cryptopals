package set1

import (
	"encoding/base64"
)

// hexToDec is used to lookup a single decimal value in hex.
var hexToDec = map[byte]byte{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'a': 10,
	'b': 11,
	'c': 12,
	'd': 13,
	'e': 14,
	'f': 15,
}

// HexDecode decodes a hexidecimal byte slice to a decimal byte slice
func HexDecode(text []byte) ([]byte, error) {
	// output is always half the size of input
	out := make([]byte, len(text)/2)

	// based on this article:
	// https://learn.sparkfun.com/tutorials/hexadecimal#converting-tofrom-decimal
	for i := 0; i < len(out); i++ {
		firstDigit := hexToDec[text[i*2]]
		secondDigit := hexToDec[text[i*2+1]]

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
