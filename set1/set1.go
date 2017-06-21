package set1

import (
	"encoding/base64"
)

// hexToDec is used to lookup a single decimal value in hex.
func hexToDec(char byte) byte {
	// subtract ASCII decimal codes to convert "0" into 0 and "a" into 10.
	switch {
	case char >= '0' && char <= '9':
		return char - 48
	case char >= 'a' && char <= 'f':
		return char - 87
	}

	return 0
}

// HexDecode decodes a hexidecimal byte slice to a decimal byte slice
func HexDecode(text []byte) ([]byte, error) {
	// output is always half the size of input
	out := make([]byte, len(text)/2)

	// based on this article:
	// https://learn.sparkfun.com/tutorials/hexadecimal#converting-tofrom-decimal
	for i := 0; i < len(out); i++ {
		firstDigit := hexToDec(text[i*2])
		secondDigit := hexToDec(text[i*2+1])

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
