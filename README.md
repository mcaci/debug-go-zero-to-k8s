# Debugging Go: from zero to Kubernetes

This repo contains a web application that takes a text and creates an image or a GIF with the input text in ASCII Art

## How to use it

There are some flags that can be used to customize the image that is produced.

This is the flag list:

- bgHex: this flag gets a string in the form of "0x" followed by 3 or 6 hexadecimal digits. The value is used to select the color of the background.
- fgHex: this flag gets a string in the form of "0x" followed by 3 or 6 hexadecimal digits. The value is used to select the color of the text in the foreground.
- h: this flag is just the height of the image. It's an integer.
- l: this flag is just the width of the image. It's an integer.
- o: name of the output image
- fontPath: path for the font to use to draw the text in the foreground. It should be a path to a valid TrueType font.
- fontSize: size of the font to use to draw the text in the foreground.
- xPtFactor: this number is a factor used to determine the width of the character box for each character. It is used to adjust the alignment of each character of the ASCII art text drawn.
- yPtFactor: this number is a factor used to determine the height of the character box for each character. It is used to adjust the alignment of each character of the ASCII art text drawn.
- figlet: name of the figlet font: figlets are fonts used to convert text into ASCII art. See https://github.com/common-nighthawk/go-figure/tree/master/fonts for the possible values and http://www.figlet.org/examples.html to see examples of what are the effects of these fonts.
- banner: if set to true it will produce a gif of a banner that shows the text sliding on the background, if false it will produce a png.
- blink: if set to true it will produce a gif of the text blinking on the background, if false it will produce a png.
- alt: if set to true it will produce a gif of the text blinking and switching colors with the background, if false it will produce a png.
- delay: used with `banner`, `blink` or `alt`, it indicates the delay between each frame of the gif.

These flags can be listed using the `--help` flag.
