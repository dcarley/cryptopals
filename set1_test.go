package cryptopals_test

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"

	. "github.com/dcarley/cryptopals"

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

			It("should decode base64 to text with newlines", func() {
				input := []byte("hello gopher")

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
				base64.StdEncoding.Encode(encoded, input)

				var encodedNewlines bytes.Buffer
				for i := 0; i < len(encoded); i += 2 {
					encodedNewlines.Write(encoded[i : i+2])
					encodedNewlines.WriteByte('\n')
				}
				encoded = encodedNewlines.Bytes()
				Expect(bytes.Count(encoded, []byte{byte('\n')})).To(BeNumerically(">", 4))

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
				xor1, err := HexDecode([]byte("1c0111001f010100061a024b53535009181c"))
				Expect(err).ToNot(HaveOccurred())
				xor2, err := HexDecode([]byte("686974207468652062756c6c277320657965"))
				Expect(err).ToNot(HaveOccurred())

				xor, err := FixedKeyXOR(xor1, xor2)
				Expect(err).ToNot(HaveOccurred())
				Expect(xor).To(Equal([]byte("the kid don't play")))

				Expect(HexEncode(xor)).To(Equal([]byte("746865206b696420646f6e277420706c6179")))
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
				xor, err := HexDecode([]byte("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"))
				Expect(err).ToNot(HaveOccurred())

				score, err := BruteForceSingleByteXOR(xor)
				Expect(err).ToNot(HaveOccurred())
				Expect(score.Text).To(Equal([]byte("Cooking MC's like a pound of bacon")))
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

	Describe("Challenge4", func() {
		It("should solve example", func() {
			file, err := os.Open("fixtures/s1c4")
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()

			highestScore := KeyScore{}
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				text, err := HexDecode(scanner.Bytes())
				Expect(err).ToNot(HaveOccurred())

				score, err := BruteForceSingleByteXOR(text)
				Expect(err).ToNot(HaveOccurred())

				if score.Score > highestScore.Score {
					highestScore = score
				}
			}

			Expect(highestScore.Text).To(Equal([]byte("Now that the party is jumping\n")))
		})
	})

	Describe("Challenge5", func() {
		It("should solve example", func() {
			output, err := RepeatingKeyXOR(
				[]byte(`Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`),
				[]byte("ICE"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(HexEncode(output)).To(Equal([]byte("0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f")))
		})
	})

	Describe("Challenge6", func() {
		Describe("BruteForceMultiByteXOR", func() {
			It("should solve example", func() {
				inFile, err := os.Open("fixtures/s1c6")
				Expect(err).ToNot(HaveOccurred())
				defer inFile.Close()

				b64, err := ioutil.ReadAll(inFile)
				Expect(err).ToNot(HaveOccurred())
				xor, err := Base64Decode(b64)
				Expect(err).ToNot(HaveOccurred())

				score, err := BruteForceMultiByteXOR(xor)
				Expect(err).ToNot(HaveOccurred())
				Expect(score.Key).To(Equal([]byte("Terminator X: Bring the noise")))

				plainFile, err := os.Open("fixtures/s1c6.plain")
				Expect(err).ToNot(HaveOccurred())
				defer plainFile.Close()

				plain, err := ioutil.ReadAll(plainFile)
				Expect(err).ToNot(HaveOccurred())
				Expect(score.Text).To(Equal(plain))
			})
		})

		Describe("HammingDistance", func() {
			It("should solve example", func() {
				distance, err := HammingDistance(
					[]byte("this is a test"),
					[]byte("wokka wokka!!!"),
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(distance).To(Equal(37))
			})

			It("should return an error if lengths don't match", func() {
				distance, err := HammingDistance([]byte("a"), []byte("bb"))
				Expect(err).To(MatchError("inputs must be same length: 1 != 2"))
				Expect(distance).To(Equal(0))
			})
		})

		DescribeTable("TransposeBlocks",
			func(size int, expected [][]byte) {
				Expect(
					TransposeBlocks([]byte("abcdefghijklmnopqsrtuvwxyz"), size),
				).To(
					Equal(expected),
				)
			},
			Entry("block size of 3", 3, [][]byte{
				[]byte("adgjmprvy"),
				[]byte("behknqtwz"),
				[]byte("cfilosux"),
			}),
			Entry("block size of 4", 4, [][]byte{
				[]byte("aeimquy"),
				[]byte("bfjnsvz"),
				[]byte("cgkorw"),
				[]byte("dhlptx"),
			}),
			Entry("block size of 5", 5, [][]byte{
				[]byte("afkpuz"),
				[]byte("bglqv"),
				[]byte("chmsw"),
				[]byte("dinrx"),
				[]byte("ejoty"),
			}),
		)
	})

	Describe("Challenge7", func() {
		Describe("DecryptAESECB", func() {
			It("should solve example", func() {
				inFile, err := os.Open("fixtures/s1c7")
				Expect(err).ToNot(HaveOccurred())
				defer inFile.Close()

				b64, err := ioutil.ReadAll(inFile)
				Expect(err).ToNot(HaveOccurred())
				out, err := Base64Decode(b64)
				Expect(err).ToNot(HaveOccurred())

				out, err = DecryptAESECB(out, []byte("YELLOW SUBMARINE"))
				Expect(err).ToNot(HaveOccurred())

				// Fixture generated using:
				// openssl enc -aes-128-ecb -d -base64 -in s1c7 -nosalt \
				//		-K "$(echo -n 'YELLOW SUBMARINE' | xxd -p)" > s1c7.plain
				plainFile, err := os.Open("fixtures/s1c7.plain")
				Expect(err).ToNot(HaveOccurred())
				defer plainFile.Close()

				plain, err := ioutil.ReadAll(plainFile)
				Expect(err).ToNot(HaveOccurred())
				Expect(out).To(Equal(plain))
			})
		})

		DescribeTable("PKCS7PaddingStrip",
			func(in, out []byte) {
				const blockSize = 16
				Expect(PKCS7PaddingStrip(in, blockSize)).To(Equal(out))
			},
			Entry("strips valid padding of 4 bytes",
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114, 4, 4, 4, 4},
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114},
			),
			Entry("strips valid padding of 3 bytes",
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114, 3, 3, 3},
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114},
			),
			Entry("strips valid padding as far as possible",
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114, 2, 2, 2},
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114, 2},
			),
			Entry("doesn't strip when there is no padding",
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114},
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114},
			),
			Entry("doesn't strip invalid padding when value doesn't match count",
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114, 3, 4, 4, 4},
				[]byte{104, 101, 108, 108, 111, 32, 103, 111, 112, 104, 101, 114, 3, 4, 4, 4},
			),
			Entry("doesn't strip invalid padding when value is greater than block size",
				[]byte{
					104, 101, 108, 108, 111, 32, 103, 111,
					112, 104, 101, 114, 20, 20, 20, 20,
					20, 20, 20, 20, 20, 20, 20, 20,
					20, 20, 20, 20, 20, 20, 20, 20,
				},
				[]byte{
					104, 101, 108, 108, 111, 32, 103, 111,
					112, 104, 101, 114, 20, 20, 20, 20,
					20, 20, 20, 20, 20, 20, 20, 20,
					20, 20, 20, 20, 20, 20, 20, 20,
				},
			),
		)
	})

	Describe("Challenge8", func() {
		Describe("DetectECB", func() {
			It("should solve example", func() {
				file, err := os.Open("fixtures/s1c8")
				Expect(err).ToNot(HaveOccurred())
				defer file.Close()

				var ecbLines [][]byte
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					hex := scanner.Bytes()
					text, err := HexDecode(hex)
					Expect(err).ToNot(HaveOccurred())

					if DetectECB(text) {
						buf := make([]byte, len(hex))
						copy(buf, hex)
						ecbLines = append(ecbLines, buf)
					}
				}

				Expect(ecbLines).To(Equal([][]byte{
					[]byte("d880619740a8a19b7840a8a31c810a3d08649af70dc06f4fd5d2d69c744cd283e2dd052f6b641dbf9d11b0348542bb5708649af70dc06f4fd5d2d69c744cd2839475c9dfdbc1d46597949d9c7e82bf5a08649af70dc06f4fd5d2d69c744cd28397a93eab8d6aecd566489154789a6b0308649af70dc06f4fd5d2d69c744cd283d403180c98c8f6db1f2a3f9c4040deb0ab51b29933f2c123c58386b06fba186a"),
				}))
			})
		})
	})
})
