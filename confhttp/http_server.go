package confhttp

import (
	"compress/gzip"
	"context"
	"github.com/liucxer/courier/httptransport/handlers"
	"net/http"
	"strconv"

	"github.com/liucxer/confmiddleware/confhttp/middlewares"
	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/httptransport"
	"github.com/liucxer/courier/ptr"
	_ "github.com/liucxer/courier/validator/strfmt"
)

type Server struct {
	Port            int                                       `env:",opt,expose"`
	OpenAPISpec     string                                    `env:",opt,copy"`
	Healthy         string                                    `env:",opt,healthCheck"`
	Debug           *bool                                     `env:""`
	ht              *httptransport.HttpTransport              `env:"-"`
	contextInjector func(ctx context.Context) context.Context `env:"-"`
}

func (s Server) WithContextInjector(contextInjector func(ctx context.Context) context.Context) *Server {
	s.contextInjector = contextInjector
	return &s
}

func (s *Server) LivenessCheck() map[string]string {
	statuses := map[string]string{}

	if s.ht != nil {
		statuses[s.ht.ServiceMeta.String()] = "ok"
	}

	return statuses
}

func (s *Server) SetDefaults() {
	if s.Port == 0 {
		s.Port = 80
	}

	if s.OpenAPISpec == "" {
		s.OpenAPISpec = "./openapi.json"
	}

	if s.Debug == nil {
		s.Debug = ptr.Bool(true)
	}

	if s.Healthy == "" {
		s.Healthy = "http://:" + strconv.FormatInt(int64(s.Port), 10) + "/"
	}
}

func (s *Server) Serve(router *courier.Router) error {
	ht := httptransport.NewHttpTransport()

	ht.Port = s.Port
	ht.SetDefaults()

	ht.Middlewares = []httptransport.HttpMiddleware{
		defaultCompress,
		middlewares.DefaultCORS(),
		middlewares.HealthCheckHandler(),
		middlewares.PProfHandler(*s.Debug),
		handlers.LogHandler(),
		NewContextInjectorMiddleware(s.contextInjector),
	}

	s.ht = ht

	return ht.Serve(router)
}

func defaultCompress(h http.Handler) http.Handler {
	return middlewares.CompressHandlerLevel(h, gzip.DefaultCompression)
}
