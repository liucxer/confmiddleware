package oauth

import (
	"errors"
	"github.com/go-courier/reflectx"
	"go/ast"
	"net/url"
	"reflect"
	"strings"
)

var (
	InvalidCallbackURI = errors.New("InvalidCallbackURI")
)

// swagger:strfmt url
type RedirectURI string

func (redirectURI RedirectURI) Parse() (*url.URL, error) {
	u, err := url.Parse(string(redirectURI))
	if err != nil {
		return nil, InvalidCallbackURI
	}
	return u, nil
}

func (redirectURI RedirectURI) CodeURL(code string, state string) (*url.URL, error) {
	u, err := redirectURI.Parse()
	if err != nil {
		return nil, err
	}

	v := struct {
		Code  string `json:"code"`
		State string `json:"state,omitempty"`
	}{}

	v.Code = code
	v.State = state

	q := url.Values{}

	marshalToQuery(reflect.Indirect(reflect.ValueOf(v)), q)

	u.RawQuery = q.Encode()
	return u, nil
}

func (redirectURI RedirectURI) TokenURL(token *Token, state string) (*url.URL, error) {
	u, err := redirectURI.Parse()
	if err != nil {
		return nil, err
	}

	q := url.Values{}

	v := struct {
		*Token
		State string `json:"state,omitempty"`
	}{}

	v.Token = token
	v.State = state

	marshalToQuery(reflect.Indirect(reflect.ValueOf(v)), q)

	u.Fragment = q.Encode()
	return u, nil
}

func marshalToQuery(rv reflect.Value, values url.Values) {
	tpe := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := tpe.Field(i)
		fieldValue := reflect.Indirect(rv.Field(i))

		if field.Anonymous {
			if fieldValue.Kind() == reflect.Struct {
				marshalToQuery(fieldValue, values)
				continue
			}
		}

		if ast.IsExported(field.Name) {
			tagJSON, exists := field.Tag.Lookup("json")
			if exists {
				name := strings.SplitN(tagJSON, ",", 2)[0]
				omitempty := strings.Index(tagJSON, "omitempty") > 0

				if name != "-" {
					fv := fieldValue.Interface()

					if omitempty && reflectx.IsEmptyValue(fv) {
						continue
					}

					text, _ := reflectx.MarshalText(fieldValue.Interface())

					values.Set(name, string(text))
				}
			}
		}
	}
}
