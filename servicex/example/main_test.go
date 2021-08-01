package example

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/liucxer/confmiddleware/servicex"
)

var (
	c = &struct {
		Name   string  `env:""`
		Client *Client `env:""`
		Server struct {
			OpenAPI string `env:",opt,copy"`
			Port    int    `env:",expose"`
		}
	}{}
)

func init() {
	c.Server.OpenAPI = "./openapi.json"
	c.Server.Port = 80

	servicex.SetServiceName("srv-example", ".")
	servicex.ConfP(c)
}

func TestConf(t *testing.T) {
	servicex.Execute(func(cmd *cobra.Command, args []string) {

	})

	spew.Dump(c)
}

type Client struct {
	Host string `env:",upstream"`
}

func (c *Client) SetDefaults() {
	if c.Host == "" {
		c.Host = "some-host"
	}
}

func (c *Client) Init() {
	fmt.Println("client init")
}
