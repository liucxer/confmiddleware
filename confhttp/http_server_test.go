package confhttp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
	"time"

	"github.com/liucxer/conflogger/v2"
	"github.com/go-courier/courier"
	"github.com/go-courier/envconf"
	"github.com/go-courier/httptransport"
	"github.com/go-courier/httptransport/httpx"
	. "github.com/onsi/gomega"
)

type GetSome struct {
	httpx.MethodGet
}

func (GetSome) Path() string {
	return "/some"
}

func (GetSome) Output(ctx context.Context) (interface{}, error) {
	html := httpx.NewHTML()

	return httpx.WithMetadata(
		httpx.Metadata("Cache-Control", "no-cache"),
	)(html), nil
}

type GetOther struct {
	httpx.MethodGet
}

func (GetOther) Path() string {
	return "/other"
}

func (GetOther) Output(ctx context.Context) (interface{}, error) {
	client := ClientEndpoint{
		Endpoint: envconf.Endpoint{
			Scheme:   "http",
			Hostname: "0.0.0.0",
			Port:     1234,
		},
	}

	client.SetDefaults()
	client.Init()

	_, _ = client.Do(ctx, NewRequest(http.MethodGet, "/some")).Into(nil)

	return nil, nil
}

func TestHttp(t *testing.T) {
	log := conflogger.Log{
		Level:        "debug",
		Format:       "text",
		Output:       conflogger.OutputAlways,
		ReportCaller: "disabled",
	}
	log.SetDefaults()
	log.Init()

	server := &Server{}
	server.SetDefaults()
	server.Port = 1234

	server2 := &Server{}
	server2.SetDefaults()
	server2.Port = 3456

	router := courier.NewRouter(httptransport.Group("/"))
	router.Register(courier.NewRouter(&GetSome{}))
	router.Register(courier.NewRouter(&GetOther{}))

	go func() {
		err := server.Serve(router)
		fmt.Println(err)

		time.Sleep(5 * time.Second)

		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)
	}()

	go func() {
		err := server2.Serve(router)
		fmt.Println(err)

		time.Sleep(5 * time.Second)

		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)
	}()

	client := ClientEndpoint{
		Endpoint: envconf.Endpoint{
			Scheme:   "http",
			Hostname: "0.0.0.0",
			Port:     uint16(server.Port),
		},
	}
	client.SetDefaults()
	client.Init()

	time.Sleep(1 * time.Second)

	printResp := func(resp *http.Response) {
		data, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(data))
	}

	t.Run("GetSome", func(t *testing.T) {
		meta, err := client.Do(context.Background(), NewRequest(http.MethodGet, "/some")).Into(nil)
		NewWithT(t).Expect(err).To(BeNil())

		NewWithT(t).Expect(http.Header(meta).Get("b3")).NotTo(BeEmpty())
	})

	t.Run("GetOther", func(t *testing.T) {
		meta, err := client.Do(context.Background(), NewRequest(http.MethodGet, "/other")).Into(nil)
		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(http.Header(meta).Get("b3")).NotTo(BeEmpty())
	})

	t.Run("Head", func(t *testing.T) {
		resp, err := http.Head(fmt.Sprintf("http://0.0.0.0:%d", server.Port))

		NewWithT(t).Expect(err).To(BeNil())
		printResp(resp)
	})

	t.Run("Options", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodOptions, fmt.Sprintf("http://0.0.0.0:%d/some", server.Port), nil)
		req.Header.Add("Origin", "http://localhost:3000")
		req.Header.Add("Access-Control-Request-Method", http.MethodPost)
		req.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
		resp, err := http.DefaultClient.Do(req)

		NewWithT(t).Expect(err).To(BeNil())

		printResp(resp)
	})

	time.Sleep(1 * time.Second)
}