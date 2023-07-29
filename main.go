package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/mcaci/debug-go-zero-to-k8s/gate"
	"github.com/mcaci/debug-go-zero-to-k8s/handler"
)

func main() {
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	http.Handle("/byBlink", gate.Get(&handler.BlinkBnY{}))
	http.Handle("/free", gate.Post(&handler.FreeStyle{}))
	http.Handle("/goCol", gate.Get(&handler.GoColorBanner{}))
	http.HandleFunc("/ping", func(http.ResponseWriter, *http.Request) {})

	log.Println("listenning on port 8080")
	log.Println("available endpoints: /byBlink, /free, /goCol, /ping")
	log.Println(http.ListenAndServe(":8080", nil))
}
