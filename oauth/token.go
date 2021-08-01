package oauth

import (
	"encoding/json"
	"github.com/go-courier/statuserror"
	"golang.org/x/oauth2"
)

func parseTokenAndError(t *oauth2.Token, err error) (*Token, error) {
	if err != nil {
		if retrieveErr, ok := err.(*oauth2.RetrieveError); ok {
			e := &statuserror.StatusErr{}
			unmarshalErr := json.Unmarshal(retrieveErr.Body, e)
			if unmarshalErr != nil {
				return nil, unmarshalErr
			}
			return nil, e
		}
		return nil, err
	}

	tok := &Token{
		Token: *t,
	}

	uid := t.Extra("uid")

	if uid != nil {
		tok.UID = uid.(string)
	}

	expiresIn := t.Extra("expires_in")
	if expiresIn != nil {
		tok.ExpiresIn = int(expiresIn.(float64))
	}

	return tok, nil
}

type TokenSource interface {
	Token() (*Token, error)
}

func wrapTokenSource(ts oauth2.TokenSource) TokenSource {
	return &tokenSourceWrapper{TokenSource: ts}
}

type tokenSourceWrapper struct {
	oauth2.TokenSource
}

func (tsw *tokenSourceWrapper) Token() (*Token, error) {
	return parseTokenAndError(tsw.TokenSource.Token())
}

type Token struct {
	oauth2.Token
	ExpiresIn int    `json:"expires_in,omitempty"`
	UID       string `json:"uid"`
}
