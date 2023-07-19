# Debugging Go: from zero to Kubernetes

This repo contains a web application that takes a text and creates an image or a GIF with the input text in ASCII Art

## How to use it

Three endpoints can be used:

1. /free: POST request transforming a text into a gif using the parameters passed in the JSON body
2. /byBlink: GET request transforming a text into a blinking gif with blue background and yellow text
3. /goCol: GET request transforming a text into a banner gif with light blue background and white text (Go's colors)

All endpoints provide the text via a query parameter named `text`, for example: `curl localhost:8080/free?text=help -d '{"gifType":"alt"}'`

The json body for the `/free` endpoint accepts the following inputs:

```go
struct {
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
```
