package openapi2word_test

import (
	"net/url"
	"testing"

	"github.com/liucxer/confmiddleware/openapi2word"
	"github.com/stretchr/testify/require"
)

func TestOpenApi_GenerateDoc(t *testing.T) {
	gen := openapi2word.NewGenerateOpenAPIDoc("access-demo", &url.URL{
		Scheme: "http",
		Path:   "srv-access-demo__aisys.rocktl.com/access-demo",
	}, 3)
	gen.Load()
	err := gen.GenerateClientOpenAPIDoc("./access-demo.docx")
	require.NoError(t, err)
}
