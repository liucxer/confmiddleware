package authorization

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthorizations(t *testing.T) {
	tt := require.New(t)

	auths := Authorizations{}

	auths.Add("Bearer", "xxxxx")
	auths.Add("WechatBearer", "yyyyy")

	t.Log(auths.String())

	tt.Equal(auths, ParseAuthorization(auths.String()))
	tt.Equal("xxxxx", auths.Get("bearer"))
	tt.Equal("yyyyy", auths.Get("WechatBearer"))
}
