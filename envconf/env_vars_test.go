package envconf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-courier/ptr"
	"github.com/go-courier/snapshotmacther"
	. "github.com/onsi/gomega"
)

type SubConfig struct {
	Key          string `env:""`
	Bool         bool
	Map          map[string]string
	Func         func() error
	defaultValue bool
}

func (c *SubConfig) SetDefaults() {
	c.defaultValue = true
}

type Config struct {
	Map       map[string]string
	Slice     []string `env:""`
	PtrString *string  `env:""`
	Host      string   `env:",upstream"`
	SubConfig
	Config SubConfig
}

func TestEnvVars(t *testing.T) {
	c := Config{}

	c.Key = "123456"
	c.PtrString = ptr.String("123456=")
	c.Slice = []string{"1", "2"}
	c.Config.Key = "key"
	c.Config.defaultValue = true
	c.defaultValue = true

	t.Run("Encoding", func(t *testing.T) {
		envVars := NewEnvVars("S")

		_ = EncodeEnvVars(envVars, &c)

		NewWithT(t).Expect(string(envVars.Bytes())).To(snapshotmacther.MatchSnapshot(".dotenv"))
	})

	t.Run("Decoding", func(t *testing.T) {
		envVars := NewEnvVars("S")

		_ = EncodeEnvVars(envVars, &c)

		envVars2 := EnvVarsFromEnviron("S", strings.Split(string(envVars.Bytes()), "\n"))

		c2 := Config{}
		err := DecodeEnvVars(envVars2, &c2)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(c).To(Equal(c2))
	})

	t.Run("Decoding By AES", func(t *testing.T) {
		key := GenKey(32)
		cryptor, _ := NewAESCryptor(key)

		valueEncryptedBy := fmt.Sprintf("aes/cfb %s", base64.RawStdEncoding.EncodeToString(key))

		os.Setenv("VALUE_ENCRYPTED_BY", valueEncryptedBy)

		envVars := NewEnvVars("S")

		_ = EncodeEnvVars(envVars, &c)

		dotEnv := envVars.DotEnv(func(envVar *EnvVar) string {
			b := bytes.NewBuffer(nil)
			w, _ := cryptor.EncryptFor(b)
			defer w.Close()
			_, _ = w.Write([]byte(envVar.Value))
			return base64.RawStdEncoding.EncodeToString(b.Bytes())
		})

		t.Log(string(dotEnv))

		envVars2 := EnvVarsFromEnviron("S", strings.Split(string(dotEnv), "\n"))

		c2 := Config{}
		err := DecodeEnvVars(envVars2, &c2)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(c).To(Equal(c2))
	})
}
