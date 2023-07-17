package main

import (
	"bufio"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/mcaci/debug-go-zero-to-k8s/application"
)

func main() {
	l := flag.Int("l", 0, "Length of the image")
	h := flag.Int("h", 0, "Height of the image")
	path := flag.String("o", "", "path of the output image/gif")
	bgColorHex := flag.String("bgHex", "0x4d3178", "Hexadecimal value for the background color")
	fgColorHex := flag.String("fgHex", "0xabc", "Hexadecimal value for the color of the text")
	fontPath := flag.String("fontPath", "fonts/Ubuntu-R.ttf", "path of the font to use")
	fontSize := flag.Float64("fontSize", 32.0, "font size of the output text in the image")
	xPtFactor := flag.Float64("xPtFactor", 0.5, "x size factor of one letter box")
	yPtFactor := flag.Float64("yPtFactor", 1.0, "y size factor of one letter box")
	figlet := flag.String("figlet", "banner", "name of the figlet font; see https://github.com/common-nighthawk/go-figure/tree/master/fonts for the values and http://www.figlet.org/examples.html for the actual effect")
	banner := flag.Bool("banner", false, "if true it's a banner gif, else it's a picture")
	blink := flag.Bool("blink", false, "if true it's a plain blinking gif, else it's a picture")
	alt := flag.Bool("alt", false, "if true it's a alternating colors blinking gif, else it's a picture")
	delay := flag.Int("delay", 0, "used with '-banner, '-blink' or '-alt', it indicates the delay between each frame of the gif")
	inFile := flag.String("inputFile", "", "Experimental: path of the file containing a list of strings for which to create the image/gif")

	flag.Parse()

	switch *inFile {
	case "":
		text := strings.Join(flag.Args(), " ")
		application.Create(*l, *h, text, *path, *bgColorHex, *fgColorHex, *figlet, *fontPath, *fontSize, *xPtFactor, *yPtFactor, *delay, *banner, *blink, *alt)
	default:

		f, _ := os.Open(*inFile)
		s := bufio.NewScanner(f)
		// s.Split(bufio.ScanLines)
		for s.Scan() {
			text := s.Text()
			application.Create(*l, *h, text, *path, *bgColorHex, *fgColorHex, *figlet, *fontPath, *fontSize, *xPtFactor, *yPtFactor, *delay, *banner, *blink, *alt)

		}
	}

	http.Handle("/hello", &application.HelloHandler{})
	http.ListenAndServe(":8080", nil)
}
