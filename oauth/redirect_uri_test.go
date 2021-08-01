package oauth

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"testing"
)

func TestRedirectURI(t *testing.T) {
	tt := require.New(t)
	u := RedirectURI("/cb")

	{
		s, err := u.CodeURL("CODE", "")
		tt.NoError(err)
		tt.Equal("/cb?code=CODE", s.String())
	}

	{
		s, err := u.CodeURL("CODE", "STATE")
		tt.NoError(err)
		tt.Equal("/cb?code=CODE&state=STATE", s.String())
	}

	{
		t := &Token{
			Token: oauth2.Token{
				TokenType:    "bearer",
				AccessToken:  "123",
				RefreshToken: "123",
			},
			ExpiresIn: 2000,
			UID:       "123",
		}

		s, err := u.TokenURL(t, "")
		tt.Nil(err)
		tt.Equal("/cb#access_token=123&expires_in=2000&refresh_token=123&token_type=bearer&uid=123", s.String())
	}
}
