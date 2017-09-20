package set1_test

import (
	. "github.com/dcarley/cryptopals/set1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"encoding/base64"
	"encoding/hex"
)

var _ = Describe("Set1", func() {
	Describe("Challenge1", func() {
		Describe("HexToBase64", func() {
			It("should convert example", func() {
				b64, err := HexToBase64([]byte("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"))
				Expect(err).ToNot(HaveOccurred())
				Expect(b64).To(Equal([]byte("SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t")))
			})
		})

		Describe("HexDecode", func() {
			It("should decode hex to decimal byte slice", func() {
				input := []byte("hello gopher")

				encoded := make([]byte, hex.EncodedLen(len(input)))
				hex.Encode(encoded, input)

				decoded, err := HexDecode(encoded)
				Expect(err).ToNot(HaveOccurred())
				Expect(decoded).To(Equal(input))
			})

			It("should handle uppercase and lowercase alphas", func() {
				decoded, err := HexDecode([]byte("6a6B6c6D"))
				Expect(err).ToNot(HaveOccurred())
				Expect(decoded).To(Equal([]byte("jklm")))
			})

			It("should return an error for invalid alphas", func() {
				decoded, err := HexDecode([]byte("6g"))
				Expect(err).To(MatchError("invalid hex character: g"))
				Expect(decoded).To(Equal([]byte{}))
			})

			It("should return an error on odd input sizes", func() {
				decoded, err := HexDecode([]byte("abc"))
				Expect(err).To(MatchError("input must be an even size"))
				Expect(decoded).To(Equal([]byte{}))
			})
		})

		Describe("HexEncode", func() {
			It("should encode decimal to hex byte slice", func() {
				input := []byte("hello gopher")

				encoded := make([]byte, hex.EncodedLen(len(input)))
				hex.Encode(encoded, input)

				Expect(HexEncode(input)).To(Equal(encoded))
			})
		})

		Describe("Base64Decode", func() {
			It("should decode base64 to text", func() {
				input := []byte("hello gopher")

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
				base64.StdEncoding.Encode(encoded, input)

				decoded, err := Base64Decode(encoded)
				Expect(err).ToNot(HaveOccurred())
				Expect(decoded).To(Equal(input))
			})

			It("should decode base64 to text with one character padding", func() {
				input := []byte("hello gophers")
				Expect(len(input) % 3).To(Equal(1))

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
				base64.StdEncoding.Encode(encoded, input)

				decoded, err := Base64Decode(encoded)
				Expect(err).ToNot(HaveOccurred())
				Expect(decoded).To(Equal(input))
			})

			It("should decode base64 to text with two character padding", func() {
				input := []byte("hello gophers!")
				Expect(len(input) % 3).To(Equal(2))

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
				base64.StdEncoding.Encode(encoded, input)

				decoded, err := Base64Decode(encoded)
				Expect(err).ToNot(HaveOccurred())
				Expect(decoded).To(Equal(input))
			})

			It("should return an error for invalid base64 characters", func() {
				decoded, err := Base64Decode([]byte("abc!def"))
				Expect(err).To(MatchError("invalid base64 character: !"))
				Expect(decoded).To(Equal([]byte{}))
			})
		})

		Describe("Base64Encode", func() {
			It("should encode text to base64", func() {
				input := []byte("hello gopher")

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
				base64.StdEncoding.Encode(encoded, input)

				Expect(Base64Encode(input)).To(Equal(encoded))
			})

			It("should encode text to base64 with one character padding", func() {
				input := []byte("hello gophers")
				Expect(len(input) % 3).To(Equal(1))

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
				base64.StdEncoding.Encode(encoded, input)

				Expect(Base64Encode(input)).To(Equal(encoded))
			})

			It("should encode text to base64 with two character padding", func() {
				input := []byte("hello gophers!")
				Expect(len(input) % 3).To(Equal(2))

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
				base64.StdEncoding.Encode(encoded, input)

				Expect(Base64Encode(input)).To(Equal(encoded))
			})
		})
	})

	Describe("Challenge2", func() {
		Describe("FixedKeyXOR", func() {
			It("should convert example", func() {
				xor, err := FixedKeyXOR(
					[]byte("1c0111001f010100061a024b53535009181c"),
					[]byte("686974207468652062756c6c277320657965"),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(xor).To(Equal([]byte("746865206b696420646f6e277420706c6179")))

				plain, err := HexDecode(xor)
				Expect(err).ToNot(HaveOccurred())
				Expect(plain).To(Equal([]byte("the kid don't play")))
			})

			It("should error on unequal lengths", func() {
				xor, err := FixedKeyXOR(
					[]byte("12345678"),
					[]byte("1234"),
				)
				Expect(err).To(MatchError("text and key must be same size: 8 != 4"))
				Expect(xor).To(Equal([]byte{}))
			})
		})
	})

	Describe("Challenge3", func() {
		Describe("BruteForceSingleByteXOR", func() {
			It("should convert example", func() {
				xor := []byte("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")

				key, err := BruteForceSingleByteXOR(xor)
				Expect(err).ToNot(HaveOccurred())

				text, err := RepeatingKeyXOR(xor, key)
				Expect(err).ToNot(HaveOccurred())

				plain, err := HexDecode(text)
				Expect(err).ToNot(HaveOccurred())
				Expect(plain).To(Equal([]byte("Cooking MC's like a pound of bacon")))
			})
		})

		DescribeTable("ScoreEnglish",
			func(text []byte, score int) {
				Expect(ScoreEnglish(text)).To(Equal(score))
			},
			Entry("repeated character", []byte("xxxxxxxxxxxxxxxxxxxxxxx"), 0),
			Entry("pwgen 23 -y", []byte("qui1Chux(euZae9Ua3pooqu"), 13),
			Entry("keyboard bashing", []byte("dgj lqn0[jr1n3ofe we[of w"), 12),
			Entry("numbers only", []byte("01234567890123456789012"), 0),
			Entry("proper English", []byte("I'm writing proper English"), 19),
			Entry("real sentence", []byte("This is a real sentence"), 22),
		)
	})
})
