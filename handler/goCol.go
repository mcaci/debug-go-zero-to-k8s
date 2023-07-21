package handler

import (
	"net/http"

	"github.com/mcaci/debug-go-zero-to-k8s/img"
)

type GoColorBanner struct{}

func (*GoColorBanner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const (
		delay      = 10
		figlet     = "doh"
		fontSize   = 32
		fontPath   = "./fonts/Ubuntu-R.ttf"
		xPtFactor  = 0.5
		yPtFactor  = 1
		bgColorHex = "0x00ADD8" // light blue
		fgColorHex = "0xFFF"    // white
	)
	defer r.Body.Close()
	text := r.URL.Query().Get("text")
	asciiArtLines := img.ToAsciiArt(text, figlet)
	l := img.Side(len(asciiArtLines[0]), 2*30, fontSize, xPtFactor, 0)
	h := img.Side(len(asciiArtLines), 2*30, fontSize, yPtFactor, 0)
	images := img.NewBanner(asciiArtLines, l, h, bgColorHex, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
	w.Header().Set("Content-Disposition", "attachment; filename="+img.OutFilename(text, "gif"))
	w.Header().Set("Content-Type", "application/octet-stream")
	err := img.WriteGif(w, images, delay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
