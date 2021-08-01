package oauth

import (
	"context"
	"github.com/go-courier/envconf"
	"github.com/go-courier/httptransport/client"
	"github.com/go-courier/httptransport/client/roundtrippers"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"time"
)

type OAuthConfig struct {
	Timeout time.Duration

	ClientID     string           `env:""`
	ClientSecret envconf.Password `env:""`
	AuthURL      string           `env:""`
	TokenURL     string           `env:""`
	RedirectURL  string           `env:""`

	Scopes []string

	oauth2.Config `env:"-"`
}

func (o *OAuthConfig) SetDefaults() {
	if o.AuthURL == "" {
		o.AuthURL = "http://srv-oauth.base.dev.rktl.work/oauth/authorize"
	}
	if o.TokenURL == "" {
		o.TokenURL = "http://srv-oauth.base.dev.rktl.work/oauth/token"
	}

	if o.Timeout == 0 {
		o.Timeout = 5 * time.Minute
	}
}

func (o *OAuthConfig) endpoint() oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:   o.AuthURL,
		TokenURL:  o.TokenURL,
		AuthStyle: oauth2.AuthStyleInHeader,
	}
}

func (o *OAuthConfig) withClient(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if o.Timeout > 0 && ctx.Value(oauth2.HTTPClient) == nil {
		return context.WithValue(ctx, oauth2.HTTPClient, client.GetShortConnClient(o.Timeout,
			roundtrippers.NewLogRoundTripper(logrus.StandardLogger()),
		))
	}

	return ctx
}

func (o *OAuthConfig) Init() {
	o.Config = oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: string(o.ClientSecret),
		Endpoint:     o.endpoint(),
		RedirectURL:  o.RedirectURL,
		Scopes:       o.Scopes,
	}
}

func (o *OAuthConfig) RedirectAuthCodeURL(authURL string) (*url.URL, error) {
	c := &http.Client{Timeout: o.Timeout}
	req, _ := http.NewRequest("GET", authURL, nil)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Request.URL, nil
}

func (o *OAuthConfig) AuthCodeURL(state string, uid string, opts ...oauth2.AuthCodeOption) string {
	return o.Config.AuthCodeURL(state, append([]oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("uid", uid),
	}, opts...)...)
}

func (o *OAuthConfig) ClientCredentialsToken(ctx context.Context, uid string) (*Token, error) {
	cc := &clientcredentials.Config{
		ClientID:     o.Config.ClientID,
		ClientSecret: o.Config.ClientSecret,
		Scopes:       o.Scopes,
		TokenURL:     o.Config.Endpoint.TokenURL,
		AuthStyle:    o.Config.Endpoint.AuthStyle,
		EndpointParams: url.Values{
			"uid": {uid},
		},
	}

	return parseTokenAndError(cc.Token(o.withClient(ctx)))
}

func (o *OAuthConfig) ValidateToken(ctx context.Context, accessToken string) (*Token, error) {
	cc := &clientcredentials.Config{
		ClientID:     o.Config.ClientID,
		ClientSecret: o.Config.ClientSecret,
		Scopes:       o.Scopes,
		TokenURL:     o.Config.Endpoint.TokenURL,
		AuthStyle:    o.Config.Endpoint.AuthStyle,
		EndpointParams: url.Values{
			"access_token": {accessToken},
		},
	}

	return parseTokenAndError(cc.Token(o.withClient(ctx)))
}

func (o *OAuthConfig) Exchange(ctx context.Context, code string) (*Token, error) {
	return parseTokenAndError(o.Config.Exchange(o.withClient(ctx), code))
}

func (o *OAuthConfig) PasswordCredentialsToken(ctx context.Context, username string, password string) (*Token, error) {
	return parseTokenAndError(o.Config.PasswordCredentialsToken(o.withClient(ctx), username, password))
}

func (o *OAuthConfig) TokenSource(ctx context.Context, token *Token) TokenSource {
	return wrapTokenSource(o.Config.TokenSource(o.withClient(ctx), &token.Token))
}
