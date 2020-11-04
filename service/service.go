package service

import (
	"fmt"
	"net/http"
	"strconv"
)

type Server interface {
	Run() error
}

type Service struct {}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("recovered from exception: %v", r)))
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
}

func (s *Service) Run() error {
	return http.ListenAndServe(":8080", PanicMiddleware(http.HandlerFunc(GenericHandler)))
}