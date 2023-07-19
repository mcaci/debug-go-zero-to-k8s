package handler

import (
	"image/gif"
	"io"
	"net/http"
	"os"

	"github.com/mcaci/debug-go-zero-to-k8s/img"
)

type BlinkBnY struct{}

func (*BlinkBnY) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const (
		delay      = 100
		figlet     = "doh"
		fontSize   = 32
		fontPath   = "./fonts/Ubuntu-R.ttf"
		xPtFactor  = 0.5
		yPtFactor  = 1
		bgColorHex = "0x053ba8" // blue
		fgColorHex = "0xded702" // yellow
	)
	defer r.Body.Close()
	text := r.URL.Query().Get("text")
	asciiArtLines := img.ToAsciiArt(text, figlet)
	l := img.Side(len(asciiArtLines[0]), 2*30, fontSize, xPtFactor, 0)
	h := img.Side(len(asciiArtLines), 2*30, fontSize, yPtFactor, 0)
	images := img.NewBlink(asciiArtLines, l, h, bgColorHex, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
	f, err := os.Create(img.OutFilename(text, "gif"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	delays := make([]int, len(images))
	for i := range delays {
		delays[i] = delay
	}
	err = gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	b, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
