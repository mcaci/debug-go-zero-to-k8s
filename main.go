package main

import (
	"net/http"

	"github.com/mcaci/debug-go-zero-to-k8s/gate"
	"github.com/mcaci/debug-go-zero-to-k8s/handler"
)

func main() {
	http.Handle("/byBlink", gate.Get(&handler.BlinkBnY{}))
	http.Handle("/free", gate.Post(&handler.FreeStyle{}))
	http.Handle("/goCol", gate.Get(&handler.GoColorBanner{}))
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	http.ListenAndServe(":8080", nil)
}
