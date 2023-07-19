package gate

import (
	"errors"
	"net/http"
)

var errEmptyRequest = errors.New("empty request")

func Get(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			http.Error(w, errEmptyRequest.Error(), http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "expecting GET request, was "+r.Method, http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Post(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil || r.Body == nil {
			http.Error(w, errEmptyRequest.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			http.Error(w, "expecting POST request, was "+r.Method, http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	})
}
