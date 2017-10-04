package set1

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"regexp"
	"sort"
)

// decToHex is used to lookup a single hex value in decimal.
var decToHex = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f',
}

// Encodes a decimal byte slice to a hexadecimal byte slice
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

// HexDecode decodes a hexadecimal byte slice to a decimal byte slice
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
	size := len(text) / 3 * 4
	// allow for padding if input is not divisible by 3
	if len(text)%3 != 0 {
		size += 4
	}
	out := make([]byte, size)

	for i := 0; i < len(text); i += 3 {
		// convert three 8bit characters into one 24bit value by:
		// - converting from byte (int8) to int32
		// - bit shifting to the left, to pad least significant bits, so that they don't overlap
		// - bit ORing to combine the values into one
		// - ignore the last two characters if we've reached the end of the input
		combined := int32(text[i]) << 16
		if i+1 < len(text) {
			combined |= int32(text[i+1]) << 8
		}
		if i+2 < len(text) {
			combined |= int32(text[i+2]) << 0
		}

		// split the 24bit value into four 6bit values by:
		// - bit shifting to the right, so the least significant 6bits are what we want
		// - bit ANDing against 6bits of binary 1s to extract the least significant 6bits
		// - looking up the appropriate base64 character for the value
		// - pad the last two characters if we've reached the end of the input
		outIndex := i / 3 * 4
		out[outIndex] = decToBase64[combined>>18&sixBitsMax]
		out[outIndex+1] = decToBase64[combined>>12&sixBitsMax]

		if i+1 < len(text) {
			out[outIndex+2] = decToBase64[combined>>6&sixBitsMax]
		} else {
			out[outIndex+2] = '='
		}

		if i+2 < len(text) {
			out[outIndex+3] = decToBase64[combined>>0&sixBitsMax]
		} else {
			out[outIndex+3] = '='
		}
	}

	return out
}

// base64ToDec is used to lookup a zero-indexed value for a base64 character
func base64ToDec(char byte) (int32, error) {
	switch {
	case char >= 'A' && char <= 'Z':
		return int32(char) - 'A', nil
	case char >= 'a' && char <= 'z':
		return int32(char) - 'a' + 26, nil
	case char >= '0' && char <= '9':
		return int32(char) - '0' + 52, nil
	case char == '+':
		return 62, nil
	case char == '/':
		return 63, nil
	case char == '=':
		return 0, nil
	}

	return 0, fmt.Errorf("invalid base64 character: %s", []byte{char})
}

// StripBytes removes all occurrences of a byte from a byte slice.
func StripBytes(s []byte, c byte) []byte {
	for {
		i := bytes.IndexByte(s, c)
		if i == -1 {
			break
		}

		s = append(s[:i], s[i+1:]...)
	}

	return s
}

// eightBitsMax is a 8bit value of all binary 1s:
//	- bin: 11111111
//	- dec: 255
//	- hex: 0xFF
const eightBitsMax = 0xFF

// Base64Decode decodes a base64 byte slice into an byte slice
func Base64Decode(text []byte) ([]byte, error) {
	text = StripBytes(text, '\n')

	// 4 bytes out for every 3 bytes in
	out := make([]byte, len(text)/4*3)

	for i := 0; i < len(text); i += 4 {
		// convert four 6bit characters into one 24bit value by:
		// - looking up the appropriate index for the base64 character
		// - bit shifting to the left, to pad least significant bits, so that they don't overlap
		// - bit ORing to combine the values into one
		var combined int32
		for offset, shift := range []uint{18, 12, 6, 0} {
			v, err := base64ToDec(text[i+offset])
			if err != nil {
				return []byte{}, err
			}

			combined |= v << shift
		}

		// split the 24bit value into three 8bit values by:
		// - bit shifting to the right, so the least significant 8bits are what we want
		// - bit ANDing against 8bits of binary 1s to extract the least significant 8bits
		// - converting to byte so that we get the appropriate ASCII code
		outIndex := i / 4 * 3
		out[outIndex] = byte(combined >> 16 & eightBitsMax)
		out[outIndex+1] = byte(combined >> 8 & eightBitsMax)
		out[outIndex+2] = byte(combined >> 0 & eightBitsMax)
	}

	// remove trailing padding
	out = bytes.TrimRight(out, "\x00")

	return out, nil
}

// HexToBase64 converts hexadecimal encoded text to base64.
func HexToBase64(text []byte) ([]byte, error) {
	raw, err := HexDecode(text)
	if err != nil {
		return []byte{}, err
	}

	return Base64Encode(raw), nil
}

// FixedKeyXOR encrypts some text against a key of the same size.
func FixedKeyXOR(text, key []byte) ([]byte, error) {
	if len(text) != len(key) {
		return []byte{}, fmt.Errorf("text and key must be same size: %d != %d", len(text), len(key))
	}

	return RepeatingKeyXOR(text, key)
}

// RepeatingKeyXOR encrypts some text against a repeating key of a smaller
// size.
func RepeatingKeyXOR(text, key []byte) ([]byte, error) {
	var keyIndex int
	xor := make([]byte, len(text))
	for i := 0; i < len(xor); i++ {
		// repeat the beginning of the key if we've reached the end
		if keyIndex >= len(key) {
			keyIndex = 0
		}

		xor[i] = text[i] ^ key[keyIndex]
		keyIndex++
	}

	return xor, nil
}

// ScoreEnglish returns a score indicating the likelihood that a string
// is comprised of English words by counting the most commonly occurring
// letters in the English language.
func ScoreEnglish(text []byte) int {
	re := regexp.MustCompile("(?i)[ETAOIN SHRDLU]")
	matches := re.FindAll(text, -1)

	return len(matches)
}

// KeyScore can be used to keep track of the most likely key.
type KeyScore struct {
	Score     int
	Key, Text []byte
}

// BruteForceSingleByteXOR finds the single byte key that some text has been
// XORed against.
func BruteForceSingleByteXOR(text []byte) (KeyScore, error) {
	var highestScore KeyScore

	// try all printable ASCII characters
	for key := byte(32); key <= byte(127); key++ {
		out, err := RepeatingKeyXOR(text, []byte{key})
		if err != nil {
			return highestScore, err
		}

		if score := ScoreEnglish(out); score > highestScore.Score {
			highestScore = KeyScore{
				Score: score,
				Key:   []byte{key},
				Text:  out,
			}
		}
	}

	return highestScore, nil
}

// HammingDistance returns the number of differences between two byte
// slices: https://en.wikipedia.org/wiki/Hamming_distance
func HammingDistance(one, two []byte) (int, error) {
	if len(one) != len(two) {
		return 0, fmt.Errorf("inputs must be same length: %d != %d", len(one), len(two))
	}

	var dist int
	for i := 0; i < len(one); i++ {
		// find bits that differ
		charXOR := one[i] ^ two[i]
		// select each bit from right to left
		for mask := 1; mask <= 128; mask *= 2 {
			if (charXOR & byte(mask)) > 0 {
				dist++
			}
		}
	}

	return dist, nil
}

// TransposeBlocks divides the input into blocks of size and returns a slice
// of byte slices, the first containining the first byte from each block,
// the second containing the second byte for each block, etc.
func TransposeBlocks(in []byte, size int) [][]byte {
	out := make([][]byte, size)
	for i := 0; i < len(in); i++ {
		out[i%size] = append(out[i%size], in[i])
	}

	return out
}

// KeySize can be used to keep track of the most likely key size.
type KeySize struct {
	Size               int
	NoramlisedDistance float64
}

// GuessXORKeySize guesses the most likely key size for some XORed text. The
// results are sorted in order of probability.
func GuessXORKeySize(text []byte) ([]int, error) {
	const (
		attempts = 10 // how many blocks to compare for each key size
		guesses  = 4  // how many guesses to return
		minSize  = 2  // minimum key size
		maxSize  = 40 // maximum key size
	)

	keySizes := make([]KeySize, maxSize+1-minSize)
	for keySize := minSize; keySize <= maxSize; keySize++ {
		var keyDistance int
		for i := 0; i <= attempts; i++ {
			var (
				lower  = keySize * i
				middle = keySize*i + keySize
				upper  = keySize*i + keySize*2
			)

			blockDistance, err := HammingDistance(text[lower:middle], text[middle:upper])
			if err != nil {
				return []int{}, err
			}

			keyDistance += blockDistance
		}

		keySizes[keySize-minSize] = KeySize{
			Size:               keySize,
			NoramlisedDistance: float64(keyDistance) / float64(keySize),
		}
	}

	sort.Slice(keySizes, func(i, j int) bool {
		return keySizes[i].NoramlisedDistance < keySizes[j].NoramlisedDistance
	})

	// lowest normalied hamming distance is most likely to be the correct size
	lowest := make([]int, guesses)
	for i := range lowest {
		lowest[i] = keySizes[i].Size
	}

	return lowest, nil
}

// BruteForceMultiByteXOR finds the multi byte key that some text has been
// XORed against.
func BruteForceMultiByteXOR(text []byte) (KeyScore, error) {
	keySizes, err := GuessXORKeySize(text)
	if err != nil {
		return KeyScore{}, err
	}

	highestScore := KeyScore{}
	for _, keySize := range keySizes {
		key := make([]byte, keySize)
		keyBlocks := TransposeBlocks(text, keySize)

		for i := 0; i < keySize; i++ {
			score, err := BruteForceSingleByteXOR(keyBlocks[i])
			if err != nil {
				return KeyScore{}, err
			}

			key[i] = score.Key[0]
		}

		out, err := RepeatingKeyXOR(text, key)
		if err != nil {
			return KeyScore{}, err
		}

		if score := ScoreEnglish(out); score > highestScore.Score {
			highestScore = KeyScore{
				Score: score,
				Key:   key,
				Text:  out,
			}
		}
	}

	return highestScore, nil
}

// StripPadding removes PKCS#7 padding.
func StripPadding(text []byte, blockSize int) []byte {
	// last byte may contain the amount of padding
	padLength := int(text[len(text)-1 : len(text)][0])
	if padLength > blockSize {
		// no padding
		return text
	}

	for i := len(text) - 1; i >= len(text)-padLength; i-- {
		if int(text[i]) != padLength {
			// padding incorrect, but don't indicate so because:
			// https://en.wikipedia.org/wiki/Padding_oracle_attack
			return text
		}
	}

	return text[:len(text)-padLength]
}

// DecryptAESECB decrypts some text that has been encrypted with AES in ECB
// mode.
func DecryptAESECB(text, key []byte) ([]byte, error) {
	ciph, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	blockSize := ciph.BlockSize()
	for i := 0; i < len(text); i += blockSize {
		// decrypt only does one block at a time.
		ciph.Decrypt(text[i:i+blockSize], text[i:i+blockSize])
	}

	return StripPadding(text, blockSize), nil
}

// DetectECB detects whether a byte slice has been encrypted in ECB mode by
// seeing if it has repeating blocks of data.
func DetectECB(text []byte) bool {
	const blockSize = 16

	blockMap := make(map[[blockSize]byte]int, len(text)/blockSize)
	for i := 0; i < len(text); i += blockSize {
		var key [blockSize]byte
		copy(key[:], text[i:i+blockSize])
		blockMap[key]++
	}

	for _, count := range blockMap {
		if count > 1 {
			return true
		}
	}

	return false
}
