package service

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type Server interface {
	Run() error
}

type Service struct {}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	return
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		wp := &responseWriter{ResponseWriter: w}
		t := time.Now()
		next.ServeHTTP(wp, r)
		elapsed := float64(time.Since(t).Nanoseconds()) / 1000000.0
		log.Info().
			Str("method", fmt.Sprintf("%v", r.Method)).
			Int("status", wp.status).
			Str("url", fmt.Sprintf("%v", r.URL)).
			Interface("ms", elapsed).
			Msg("logger")
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Str("debug", string(debug.Stack())).
					Msg("recovered from exception")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	a := []string{"a", "b"}
	val, err := strconv.Atoi(r.FormValue("i"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if _, err := w.Write([]byte(a[val])); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Service) Run() error {
	return http.ListenAndServe(":8080", PanicMiddleware(LogMiddleware(http.HandlerFunc(GenericHandler))))
}