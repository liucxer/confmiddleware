package confhttp

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/liucxer/courier/httptransport/httpx"
	"github.com/liucxer/courier/metax"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/api/trace"
)

func NewLogRoundTripper(logger *logrus.Entry) func(roundTripper http.RoundTripper) http.RoundTripper {
	return func(roundTripper http.RoundTripper) http.RoundTripper {
		return &LogRoundTripper{
			logger:           logger,
			nextRoundTripper: roundTripper,
		}
	}
}

type LogRoundTripper struct {
	logger           *logrus.Entry
	nextRoundTripper http.RoundTripper
}

func (rt *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	startedAt := time.Now()

	ctx := req.Context()

	// inject h3 form context
	(&b3.B3{InjectEncoding: b3.B3SingleHeader}).Inject(ctx, req.Header)

	resp, err := rt.nextRoundTripper.RoundTrip(req)

	level, _ := logrus.ParseLevel(strings.ToLower(req.Header.Get("x-log-level")))
	if level == logrus.PanicLevel {
		level = rt.logger.Logger.Level
	}

	cost := time.Since(startedAt)
	if err == nil {
		// extract b3 to ctx
		ctx = (&b3.B3{}).Extract(ctx, resp.Header)
	}

	logger := rt.logger.WithContext(ctx).WithFields(logrus.Fields{
		"cost":   fmt.Sprintf("%0.3fms", float64(cost/time.Millisecond)),
		"method": req.Method,
		"url":    omitAuthorization(req.URL),
	})

	if err == nil {
		if level >= logrus.InfoLevel {
			logger.Infof("success")
		}
	} else {
		if level >= logrus.WarnLevel {
			logger.Warnf("do http request failed %s", err)
		}
	}

	return resp, err
}

func LogHandler(logger *logrus.Entry, tracer trace.Tracer) func(handler http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			ctx = (&b3.B3{}).Extract(ctx, req.Header)
			startAt := time.Now()

			ctx, span := tracer.Start(ctx, "UnknownOperation", trace.WithTimestamp(startAt))
			defer func() {
				span.End(trace.WithTimestamp(time.Now()))
			}()

			loggerRw := newLoggerResponseWriter(rw)

			// for async pick
			(&b3.B3{InjectEncoding: b3.B3SingleHeader}).Inject(ctx, loggerRw.Header())

			ctx = metax.ContextWithMeta(ctx, metax.ParseMeta(span.SpanContext().TraceID.String()))

			nextHandler.ServeHTTP(loggerRw, req.WithContext(ctx))

			operator := metax.ParseMeta(loggerRw.Header().Get("X-Meta")).Get("operator")
			if operator != "" {
				span.SetName(operator)
			}

			level, _ := logrus.ParseLevel(strings.ToLower(req.Header.Get("x-log-level")))
			if level == logrus.PanicLevel {
				level = logger.Logger.Level
			}

			duration := time.Since(startAt)

			logger := logger.WithContext(metax.ContextWithMeta(ctx, metax.ParseMeta(loggerRw.Header().Get("X-Meta"))))

			header := req.Header

			fields := logrus.Fields{
				"tag":         "access",
				"remote_ip":   httpx.ClientIP(req),
				"cost":        fmt.Sprintf("%0.3fms", float64(duration/time.Millisecond)),
				"method":      req.Method,
				"request_uri": omitAuthorization(req.URL),
				"user_agent":  header.Get(httpx.HeaderUserAgent),
			}

			fields["status"] = loggerRw.statusCode

			if loggerRw.errMsg.Len() > 0 {
				if loggerRw.statusCode >= http.StatusInternalServerError {
					if level >= logrus.ErrorLevel {
						logger.WithFields(fields).Error(loggerRw.errMsg.String())
					}
				} else {
					if level >= logrus.WarnLevel {
						logger.WithFields(fields).Warn(loggerRw.errMsg.String())
					}
				}
			} else {
				if level >= logrus.InfoLevel {
					logger.WithFields(fields).Info()
				}
			}
		})
	}
}

func newLoggerResponseWriter(rw http.ResponseWriter) *loggerResponseWriter {
	h, hok := rw.(http.Hijacker)
	if !hok {
		h = nil
	}

	f, fok := rw.(http.Flusher)
	if !fok {
		f = nil
	}

	return &loggerResponseWriter{
		ResponseWriter: rw,
		Hijacker:       h,
		Flusher:        f,
	}
}

type loggerResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	http.Flusher

	headerWritten bool
	statusCode    int
	errMsg        bytes.Buffer
}

func (rw *loggerResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *loggerResponseWriter) WriteHeader(statusCode int) {
	rw.writeHeader(statusCode)
}

func (rw *loggerResponseWriter) Write(data []byte) (int, error) {
	if rw.statusCode >= http.StatusBadRequest {
		rw.errMsg.Write(data)
	}
	return rw.ResponseWriter.Write(data)
}

func (rw *loggerResponseWriter) writeHeader(statusCode int) {
	if !rw.headerWritten {
		rw.ResponseWriter.WriteHeader(statusCode)
		rw.statusCode = statusCode
		rw.headerWritten = true
	}
}

func omitAuthorization(u *url.URL) string {
	query := u.Query()
	query.Del("authorization")
	u.RawQuery = query.Encode()
	return u.String()
}
