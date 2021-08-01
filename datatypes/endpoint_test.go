package datatypes

import (
	"net/url"
	"testing"

	. "github.com/onsi/gomega"
)

func TestEndpoint(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		var id = Endpoint{
			Scheme:   "stmps",
			Hostname: "smtp.exmail.qq.com",
			Port:     465,
		}
		NewWithT(t).Expect(id.String()).To(Equal("stmps://smtp.exmail.qq.com:465"))
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		var id Endpoint

		err := id.UnmarshalText([]byte("stmps://smtp.exmail.qq.com:465"))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(id).To(Equal(Endpoint{
			Scheme:   "stmps",
			Hostname: "smtp.exmail.qq.com",
			Port:     465,
		}))
	})

	t.Run("UnmarshalText full", func(t *testing.T) {
		var id Endpoint

		err := id.UnmarshalText([]byte("postgres://username:Pass1234-_.~!*'();&@127.0.0.1:5432/base_name/xxx?sslmode=disable"))
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(id).To(Equal(Endpoint{
			Scheme:   "postgres",
			Hostname: "127.0.0.1",
			Password: "Pass1234-_.~!*'();&",
			Username: "username",
			Port:     5432,
			Base:     "base_name",
			Extra:    url.Values{"sslmode": {"disable"}},
		}))

		NewWithT(t).Expect(id.String()).To(Equal("postgres://username:Pass1234-_.~!*'();&@127.0.0.1:5432/base_name?sslmode=disable"))
		NewWithT(t).Expect(id.SecurityString()).To(Equal("postgres://username:------@127.0.0.1:5432/base_name?sslmode=disable"))
	})
}
