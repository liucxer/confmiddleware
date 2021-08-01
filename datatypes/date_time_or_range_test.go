package datatypes

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestDateTimeOrRange(t *testing.T) {
	t.Run("exactly value", func(t *testing.T) {
		tr := DateTimeOrRange{}

		err := tr.UnmarshalText([]byte("2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(tr.From.String()).To(Equal("2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(tr.To).To(BeZero())
		NewWithT(t).Expect(tr.ExclusiveFrom).To(BeFalse())
		NewWithT(t).Expect(tr.ExclusiveTo).To(BeFalse())
		NewWithT(t).Expect(tr.Exactly).To(BeTrue())

		s, err := tr.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(s)).To(Equal("2019-12-04T00:00:00+08:00"))
	})

	t.Run("from", func(t *testing.T) {
		tr := DateTimeOrRange{}

		err := tr.UnmarshalText([]byte("2019-12-04T00:00:00+08:00.."))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(tr.From.String()).To(Equal("2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(tr.To).To(BeZero())

		NewWithT(t).Expect(tr.ExclusiveFrom).To(BeFalse())
		NewWithT(t).Expect(tr.ExclusiveTo).To(BeFalse())
		NewWithT(t).Expect(tr.Exactly).To(BeFalse())

		s, err := tr.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(s)).To(Equal("2019-12-04T00:00:00+08:00.."))
	})

	t.Run("exclusive from", func(t *testing.T) {
		tr := DateTimeOrRange{}

		err := tr.UnmarshalText([]byte("2019-12-04T00:00:00+08:00<.."))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(tr.From.String()).To(Equal("2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(tr.ExclusiveFrom).To(BeTrue())
		NewWithT(t).Expect(tr.To).To(BeZero())
		NewWithT(t).Expect(tr.ExclusiveTo).To(BeFalse())

		s, err := tr.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(s)).To(Equal("2019-12-04T00:00:00+08:00<.."))
	})

	t.Run("to", func(t *testing.T) {
		tr := DateTimeOrRange{}

		err := tr.UnmarshalText([]byte("..2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(tr.From).To(BeZero())
		NewWithT(t).Expect(tr.To.String()).To(Equal("2019-12-04T00:00:00+08:00"))

		s, err := tr.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(s)).To(Equal("..2019-12-04T00:00:00+08:00"))
	})

	t.Run("full range", func(t *testing.T) {
		tr := DateTimeOrRange{}

		err := tr.UnmarshalText([]byte("2019-12-04T00:00:00+08:00<..<2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(tr.From.String()).To(Equal("2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(tr.ExclusiveFrom).To(BeTrue())
		NewWithT(t).Expect(tr.To.String()).To(Equal("2019-12-04T00:00:00+08:00"))
		NewWithT(t).Expect(tr.ExclusiveTo).To(BeTrue())

		s, err := tr.MarshalText()
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(string(s)).To(Equal("2019-12-04T00:00:00+08:00<..<2019-12-04T00:00:00+08:00"))
	})
}
