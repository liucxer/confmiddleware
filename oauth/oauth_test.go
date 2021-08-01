package oauth

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

var uid = "1095624675796455428"

var oauth = &OAuthConfig{
	ClientID:     "1098508120092577795",
	ClientSecret: "V49bfu9NZ9dUU1gBc17c8hqACVWiqTSH",
	RedirectURL:  "http://localhost:8000/callback",
	Scopes: []string{
		"all",
	},
}

func init() {
	oauth.SetDefaults()
	oauth.Init()
}

func TestOAuthFlow(t *testing.T) {
	ctx := context.Background()
	tt := require.New(t)

	authURL := oauth.AuthCodeURL("TEST_MOCK", uid)

	fmt.Println(authURL)

	u, err := oauth.RedirectAuthCodeURL(authURL)
	tt.NoError(err)

	fmt.Println(u.String())

	query := u.Query()
	code := query.Get("code")

	token, err := oauth.Exchange(ctx, code)
	tt.NoError(err)

	token2, err := oauth.ValidateToken(ctx, token.AccessToken)
	tt.NoError(err)

	tt.Equal(token2.UID, token.UID)

	refreshedToken, err := oauth.TokenSource(ctx, token).Token()
	tt.NoError(err)
	tt.Equal(token, refreshedToken)
}

func TestOAuthFlowClientCredentialsToken(t *testing.T) {
	tt := require.New(t)

	token, err := oauth.ClientCredentialsToken(context.Background(), uid)
	tt.NoError(err)

	spew.Dump(token)
}

func TestOAuthError(t *testing.T) {
	tt := require.New(t)

	_, err := oauth.ValidateToken(context.Background(), "~~")
	tt.Error(err)

	spew.Dump(err)
}
