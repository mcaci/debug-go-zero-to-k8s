package handler

import (
	"encoding/json"
	"image"
	"io"
	"net/http"

	"github.com/mcaci/debug-go-zero-to-k8s/img"
)

func FreeStyle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var jb struct {
		Delay      int     `json:"delay"`
		Figlet     string  `json:"figlet"`
		FontSize   float64 `json:"fontSize"`
		FontPath   string  `json:"fontPath"`
		GifType    string  `json:"gifType"`
		XPtFactor  float64 `json:"xPtF"`
		YPtFactor  float64 `json:"yPtF"`
		BgColorHex string  `json:"bgHex"`
		FgColorHex string  `json:"fgHex"`
	}
	err = json.Unmarshal(body, &jb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defaultIfNone(&jb)
	defer r.Body.Close()
	text := r.URL.Query().Get("text")
	asciiArtLines := img.ToAsciiArt(text, jb.Figlet)
	l := img.Side(len(asciiArtLines[0]), 2*30, jb.FontSize, jb.XPtFactor, 0)
	h := img.Side(len(asciiArtLines), 2*30, jb.FontSize, jb.YPtFactor, 0)
	var images []*image.Paletted
	switch jb.GifType {
	case "alt":
		images = img.NewAlt(asciiArtLines, l, h, jb.BgColorHex, jb.FgColorHex, jb.FontPath, jb.FontSize, jb.XPtFactor, jb.YPtFactor)
	case "banner":
		images = img.NewBlink(asciiArtLines, l, h, jb.BgColorHex, jb.FgColorHex, jb.FontPath, jb.FontSize, jb.XPtFactor, jb.YPtFactor)
	case "blink":
		images = img.NewBanner(asciiArtLines, l, h, jb.BgColorHex, jb.FgColorHex, jb.FontPath, jb.FontSize, jb.XPtFactor, jb.YPtFactor)
	default:
		http.Error(w, "unsupported gif type "+jb.GifType, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+img.OutFilename(text, "gif"))
	w.Header().Set("Content-Type", "application/octet-stream")
	err = img.WriteGif(w, images, jb.Delay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func defaultIfNone(data *struct {
	Delay      int     `json:"delay"`
	Figlet     string  `json:"figlet"`
	FontSize   float64 `json:"fontSize"`
	FontPath   string  `json:"fontPath"`
	GifType    string  `json:"gifType"`
	XPtFactor  float64 `json:"xPtF"`
	YPtFactor  float64 `json:"yPtF"`
	BgColorHex string  `json:"bgHex"`
	FgColorHex string  `json:"fgHex"`
}) {
	const (
		delay      = 100
		figlet     = "doh"
		fontSize   = 32
		fontPath   = "./fonts/Ubuntu-R.ttf"
		gifType    = "blink"
		xPtFactor  = 0.5
		yPtFactor  = 1
		bgColorHex = "0x00ADD8" // light blue
		fgColorHex = "0xFFF"    // white
	)
	if data.Delay == 0 {
		data.Delay = delay
	}
	if data.Figlet == "" {
		data.Figlet = figlet
	}
	if data.FontSize == 0 {
		data.FontSize = fontSize
	}
	if data.FontPath == "" {
		data.FontPath = fontPath
	}
	if data.GifType == "" {
		data.GifType = gifType
	}
	if data.XPtFactor == 0 {
		data.XPtFactor = xPtFactor
	}
	if data.YPtFactor == 0 {
		data.YPtFactor = yPtFactor
	}
	if data.BgColorHex == "" {
		data.BgColorHex = bgColorHex
	}
	if data.FgColorHex == "" {
		data.FgColorHex = fgColorHex
	}
}
