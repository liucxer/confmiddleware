package datatypes

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBinary(t *testing.T) {
	b := Binary("aaaaaaaaa")

	bytes, _ := b.MarshalText()
	NewWithT(t).Expect(string(bytes)).To(Equal(string(b)))
}
