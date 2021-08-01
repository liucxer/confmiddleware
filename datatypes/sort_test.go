package datatypes

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestSort(t *testing.T) {
	sort := Sort{}

	err := sort.UnmarshalText([]byte("createdAt!asc"))
	NewWithT(t).Expect(err).To(BeNil())

	NewWithT(t).Expect(sort.By).To(Equal("createdAt"))
	NewWithT(t).Expect(sort.Asc).To(BeTrue())

	s, err := sort.MarshalText()
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(string(s)).To(Equal("createdAt!asc"))
}
