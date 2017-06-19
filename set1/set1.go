package set1

import (
	"encoding/base64"
	"encoding/hex"
)

// HexToBase64 converts hexidecimal encoded text to base64.
func HexToBase64(text []byte) ([]byte, error) {
	raw := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(raw, text)
	if err != nil {
		return []byte{}, err
	}

	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	base64.StdEncoding.Encode(b64, raw)

	return b64, nil
}
