package set1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSet1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Set1 Suite")
}
