package cryptopals_test

import (
	. "github.com/dcarley/cryptopals"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set2", func() {
	Describe("Challenge9", func() {
		DescribeTable("PKCS7Padding",
			func(in []byte, blockSize int, out []byte) {
				Expect(PKCS7Padding(in, blockSize)).To(Equal(out))
			},
			Entry("16 byte text and 20 byte block",
				[]byte("YELLOW SUBMARINE"), 20,
				[]byte("YELLOW SUBMARINE\x04\x04\x04\x04"),
			),
			Entry("16 byte text and 24 byte block",
				[]byte("YELLOW SUBMARINE"), 24,
				[]byte("YELLOW SUBMARINE\x08\x08\x08\x08\x08\x08\x08\x08"),
			),
			Entry("16 byte text and 6 byte block",
				[]byte("YELLOW SUBMARINE"), 6,
				[]byte("YELLOW SUBMARINE\x02\x02"),
			),
			Entry("15 byte text and 8 byte block",
				[]byte("YELLOW SUBMARIN"), 8,
				[]byte("YELLOW SUBMARIN\x01"),
			),
			Entry("16 byte text and 16 byte block",
				[]byte("YELLOW SUBMARINE"), 16,
				[]byte("YELLOW SUBMARINE"),
			),
			Entry("16 byte text and 8 byte block",
				[]byte("YELLOW SUBMARINE"), 8,
				[]byte("YELLOW SUBMARINE"),
			),
		)
	})
})
