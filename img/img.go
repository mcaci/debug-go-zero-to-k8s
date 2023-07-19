package img

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Create(l, h int, text, bgColorHex, fgColorHex, figlet, fontPath string, fontSize, xPtFactor, yPtFactor float64, delay int, banner, blink, alt bool) {
	asciiArtLines := ToAsciiArt(text, figlet)
	lImg := Side(len(asciiArtLines[0]), 2*30, fontSize, xPtFactor, l)
	hImg := Side(len(asciiArtLines), 2*30, fontSize, yPtFactor, h)
	switch {
	case banner:
		images := NewBanner(asciiArtLines, lImg, hImg, bgColorHex, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
		writeGif(images, OutFilename(text, "gif"), delay, 10)
	case blink:
		images := NewBlink(asciiArtLines, lImg, hImg, bgColorHex, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
		writeGif(images, OutFilename(text, "gif"), delay, 75)
	case alt:
		images := NewAlt(asciiArtLines, lImg, hImg, bgColorHex, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
		writeGif(images, OutFilename(text, "gif"), delay, 100)
	default:
		image := NewPng(asciiArtLines, lImg, hImg, bgColorHex, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
		writePng(image, OutFilename(text, "png"))
	}
}

func OutFilename(text, ext string) string {
	text = cases.Title(language.English).String(text)
	return "out/" + strings.Replace(text, " ", "", -1) + "." + ext
}

func NewBanner(asciiArtLines []string, l, h int, bgColorHex, fgColorHex, fontPath string, fontSize, xPtFactor, yPtFactor float64) []*image.Paletted {
	d := int(float64(l) / fontSize)
	nFrames := len(asciiArtLines[0])
	var images []*image.Paletted
	for i := 0; i < nFrames; i += 2 {
		img, err := setupBG(bgColorHex, l/2, h)
		if err != nil {
			log.Fatal(err)
		}
		err = drawFG(asciiArtLines, i, i+d, img, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, img)
	}
	return images
}

func NewBlink(asciiArtLines []string, l, h int, bgColorHex, fgColorHex, fontPath string, fontSize, xPtFactor, yPtFactor float64) []*image.Paletted {
	const nFrames = 2
	var images []*image.Paletted
	for i := 0; i < nFrames; i++ {
		img, err := setupBG(bgColorHex, l, h)
		if err != nil {
			log.Fatal(err)
		}
		switch i % 2 {
		case 0:
			err = drawFG(asciiArtLines, 0, 0, img, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
			if err != nil {
				log.Fatal(err)
			}
		default:
			// do nothing (just background)
		}
		images = append(images, img)
	}
	return images
}

func NewAlt(asciiArtLines []string, l, h int, bgColorHex, fgColorHex, fontPath string, fontSize, xPtFactor, yPtFactor float64) []*image.Paletted {
	const nFrames = 2
	var images []*image.Paletted
	for i := 0; i < nFrames; i++ {
		var bgColor0x, fgColor0x string
		switch i % 2 {
		case 0:
			bgColor0x, fgColor0x = bgColorHex, fgColorHex // same as params
		default:
			bgColor0x, fgColor0x = fgColorHex, bgColorHex // switch back and front colors
		}
		img, err := setupBG(bgColor0x, l, h)
		if err != nil {
			log.Fatal(err)
		}
		err = drawFG(asciiArtLines, 0, 0, img, fgColor0x, fontPath, fontSize, xPtFactor, yPtFactor)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, img)
	}
	return images
}

func NewPng(asciiArtLines []string, l, h int, bgColorHex, fgColorHex, fontPath string, fontSize, xPtFactor, yPtFactor float64) *image.Paletted {
	img, err := setupBG(bgColorHex, l, h)
	if err != nil {
		log.Fatal(err)
	}
	err = drawFG(asciiArtLines, 0, 0, img, fgColorHex, fontPath, fontSize, xPtFactor, yPtFactor)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func writeGif(images []*image.Paletted, path string, delay, defaultDelay int) {
	if delay == 0 {
		delay = defaultDelay
	}
	f := mustFile(path)
	defer f.Close()
	delays := make([]int, len(images))
	for i := range delays {
		delays[i] = delay
	}
	err := gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func writePng(image *image.Paletted, path string) {
	f := mustFile(path)
	defer f.Close()
	err := png.Encode(f, image)
	if err != nil {
		log.Fatal(err)
	}
}

func mustFile(name string) *os.File {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func ToAsciiArt(text, figlet string) []string {
	fig := figure.NewFigure(text, figlet, true)
	asciiArtLines := fig.Slicify()
	var maxLineLen int
	for i := range asciiArtLines {
		if maxLineLen >= len(asciiArtLines[i]) {
			continue
		}
		maxLineLen = len(asciiArtLines[i])
	}
	for i := range asciiArtLines {
		asciiArtLines[i] += strings.Repeat(" ", maxLineLen-len(asciiArtLines[i]))
	}
	return asciiArtLines
}

func Side(n, offset int, fontSize, ptFactor float64, defaultSide int) int {
	if defaultSide != 0 {
		return defaultSide
	}
	return n*int(fontSize*ptFactor) + offset
}

func setupBG(bgHex string, l, h int) (*image.Paletted, error) {
	c, err := parseHexColor(bgHex)
	if err != nil {
		return nil, err
	}
	bg := image.NewPaletted(image.Rect(0, 0, l, h), palette.Plan9)
	draw.Draw(bg, bg.Bounds(), image.NewUniform(c), image.Pt(0, 0), draw.Src)
	return bg, nil
}

func drawFG(lines []string, s, e int, bg draw.Image, fgHex, fontPath string, fontSize, xPtFactor, yPtFactor float64) error {
	c, err := fgContext(bg, fgHex, fontPath, fontSize)
	if err != nil {
		return err
	}
	textXOffset := 10
	textYOffset := 30 + int(c.PointToFixed(fontSize)>>6) // Note shift/truncate 6 bits first

	pt := freetype.Pt(textXOffset, textYOffset)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		switch {
		case s < e && e < len(line):
			line = line[s:e]
		case s < e && e >= len(line):
			line = line[s:]
		}
		// log.Print(line)
		startX := pt.X
		for _, char := range line {
			_, err := c.DrawString(string(char), pt)
			if err != nil {
				return err
			}
			pt.X += c.PointToFixed(fontSize * xPtFactor)
		}
		pt.X = startX
		pt.Y += c.PointToFixed(fontSize * yPtFactor)
	}
	return nil
}

func fgContext(bg draw.Image, fgColorHex, fontPath string, fontSize float64) (*freetype.Context, error) {
	c := freetype.NewContext()
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	c.SetFont(f)
	c.SetDPI(72)
	c.SetFontSize(fontSize)
	c.SetClip(bg.Bounds())
	c.SetDst(bg)
	fgColor, err := parseHexColor(fgColorHex)
	if err != nil {
		return nil, err
	}
	c.SetSrc(image.NewUniform(fgColor))
	c.SetHinting(font.HintingNone)
	return c, nil
}

func parseHexColor(hex string) (color.RGBA, error) {
	var c color.RGBA
	var err error
	c.A = 0xff
	switch len(hex) {
	case 8:
		_, err = fmt.Sscanf(hex, "0x%02x%02x%02x", &c.R, &c.G, &c.B)
		return c, err
	case 5:
		_, err = fmt.Sscanf(hex, "0x%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
		return c, err
	default:
		return color.RGBA{}, fmt.Errorf("invalid length, must be 8 or 5")
	}
}
