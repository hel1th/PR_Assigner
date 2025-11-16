package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func RegisterMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Use(middleware.RequestLogger(&zerologFormatter{}))
}

type zerologFormatter struct{}

func (f *zerologFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &zerologLogEntry{req: r}
}

type zerologLogEntry struct {
	req *http.Request
}

func (e *zerologLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	log.Info().
		Str("method", e.req.Method).
		Str("url", e.req.URL.String()).
		Int("status", status).
		Int("bytes", bytes).
		Dur("duration", elapsed).
		Msg("request")
}

func (e *zerologLogEntry) Panic(v interface{}, stack []byte) {
	log.Error().
		Interface("panic", v).
		Bytes("stack", stack).
		Msg("panic")
}
