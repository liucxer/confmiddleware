package main

import (
	"fmt"

	"github.com/liucxer/confmiddleware/scaffold/pkg/appinfo"
)

var (
	client = &Client{}
	server = &Server{
		OpenAPI: "./openapi.json",
		Port:    80,
	}

	app = appinfo.New(
		appinfo.WithMainRoot("./"),
		appinfo.WithVersion("0.0.0"),
	)
)

func init() {
	app.ConfP(
		client,
		server,
	)

	app.AddCommand("migrate", func(args ...string) {
		fmt.Println("migrate")
	})
}

func main() {
	app.Execute(func(args ...string) {
		fmt.Println("main", app)
	})
}

type Server struct {
	OpenAPI string `env:",opt,copy"`
	Port    int    `env:",expose"`
}

type Client struct {
	Endpoint string `env:""`
}

func (c *Client) SetDefaults() {
	if c.Endpoint == "" {
		c.Endpoint = "http://localhost"
	}
}

func (c *Client) Init() {
	fmt.Println("client init")
}
