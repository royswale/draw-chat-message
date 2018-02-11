package main

import (
	"os"
	"log"
	"bufio"
	"flag"
	"strings"

	"image"
	"image/draw"
	"image/jpeg"
	"image/color"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font/gofont/gobold"
	"github.com/golang/freetype/truetype"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	// chat message image size
	imageWidth = 500
	imageHeight = 200

	// offset of the ZP
	offsetX = 20
	offsetY = 10
)

var (
	nameColor = color.RGBA{0, 119, 170, 255} // #07a
	dateColor = color.RGBA{113, 140, 147, 255} // #858C93
)

var dpi = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")

func main() {
	flag.Parse()

	// Parse all font data.
	fontBold, err := truetype.Parse(gobold.TTF)
	if err != nil {
		log.Fatalf("parse font bytes failed: %v", err)
	}
	fontRegular, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("parse font bytes failed: %v", err)
	}

	// New a RGBA image with the defined size.
	rgbaImage := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	backgroundImage := image.White
	draw.Draw(rgbaImage, rgbaImage.Bounds(), backgroundImage, image.ZP, draw.Src)

	// Draw the avatar image.
	avatarImageReader, err := os.Open("testpic/lord_40.jpg")
	if err != nil {
		log.Fatalf("open avatar image failed: %v", err)
	}
	defer avatarImageReader.Close()

	avatarImage, _, err := image.Decode(avatarImageReader)
	if err != nil {
		log.Fatalf("decode avatar image failed: %v", err)
	}

	draw.Draw(rgbaImage, rgbaImage.Bounds(), avatarImage,
		rgbaImage.Bounds().Min.Sub(image.Pt(offsetX, offsetY)), draw.Src)

	// --- Draw the name
	nameDrawer := newDrawer(rgbaImage, fontBold, nameColor, 18)
	nameX := offsetX + 50 + 10 // image width is 50px and offset is 10px
	nameY := offsetY + 18      // font size is 18px
	nameDrawer.Dot = fixed.P(nameX, nameY)
	nameDrawer.DrawString("Jihan Wu")

	// --- Draw the date
	dateDrawer := newDrawer(rgbaImage, fontRegular, dateColor, 18)
	dateDrawer.Dot = fixed.Point26_6{
		// offset is 10px
		X: fixed.I(nameX + 10) + nameDrawer.MeasureString("Jihan Wu"),
		// font size is 18px
		Y: fixed.I(offsetY + 18),
	}
	dateDrawer.DrawString("1/24/17")

	// --- Draw the message
	msg := "Don't play hatred\nMake BCH better"
	messageDrawer := newDrawer(rgbaImage, fontBold, color.Black, 48)
	for i, line := range strings.Split(msg, "\n") {
		// font size is 48px, spacing of lines is 5px
		messageDrawer.Dot = fixed.P(nameX, nameY + (48 + 5) * (i+1))
		messageDrawer.DrawString(line)
	}

	// Save image to disk.
	outFile, err := os.Create("out.jpeg")
	if err != nil {
		log.Fatalf("creat file failed: %v", err)
	}
	defer outFile.Close()

	bufferWriter := bufio.NewWriter(outFile)
	err = jpeg.Encode(bufferWriter, rgbaImage, nil)
	if err != nil {
		log.Fatalf("encode image failed: %v", err)
	}
	err = bufferWriter.Flush()
	if err != nil {
		log.Fatalf("flush buffer to disk failed: %v", err)
	}

	err = open.Run("out.jpeg")
	if err != nil {
		log.Fatalf("open image failed: %v", err)
	}
}

// New a text drawer, return a font.Drawer
// - `dstImage` is the dst image drawing on
// - new a font.Face with `f`
// - `fontSize` is the amount font size in pixels
func newDrawer(dstImage *image.RGBA, f *truetype.Font, fontColor color.Color, fontSize float64) *font.Drawer {
	return &font.Drawer{
		Dst: dstImage,
		Src: image.NewUniform(fontColor),
		Face: truetype.NewFace(f, &truetype.Options{
			Size: pixelsToPoints(fontSize),
			DPI: *dpi,
			Hinting: font.HintingFull,
		}),
	}
}

// Convert pixels to points (font size uint)
// - `pt` / 72 * `DPI` = `px`
func pixelsToPoints(px float64) float64{
	return px * 72 / *dpi
}
