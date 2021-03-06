package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"image/jpeg"

	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"github.com/skratchdot/open-golang/open"
)

const (
	// chat message image size
	imageWidth  = 500
	imageHeight = 200

	// offset of the ZP
	offsetX = 20
	offsetY = 10
)

var (
	nameColor = color.RGBA{0, 119, 170, 255}   // hex value: #07a
	dateColor = color.RGBA{113, 140, 147, 255} // hex value: #858C93
)

var (
	imageSource = flag.String("image", "", "avatar image source: file | url")
	name        = flag.String("name", "Kevin", "message user name")
	date        = flag.String("date", "1/24/2017", "message date")
	content     = flag.String("content", "", "message content")
	outputName  = flag.String("output", "out", "output file name (with no suffix)")
	isOpen      = flag.Bool("open", true, "open file with the default tool")
	stdout      = flag.Bool("stdout", false, "write output to stdout")
	dpi         = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
)

func main() {
	flag.Parse()

	// Parse all font data.
	fontBold, err := truetype.Parse(gobold.TTF)
	if err != nil {
		errorReturn("parse font bytes failed: %v", err)
	}
	fontRegular, err := truetype.Parse(goregular.TTF)
	if err != nil {
		errorReturn("parse font bytes failed: %v", err)
	}

	// New white board with the defined size.
	rgbaImage := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(rgbaImage, rgbaImage.Bounds(), image.White, image.ZP, draw.Src)

	// Draw the avatar image.
	avatarImageReader := makeImageReader()
	defer avatarImageReader.Close()

	avatarImage, _, err := image.Decode(avatarImageReader)
	if err != nil {
		errorReturn("decode avatar image failed: %v", err)
	}

	// Resize avatar image to 40 * 40 size.
	rszAvatarImage := resize.Resize(40, 40, avatarImage, resize.NearestNeighbor)

	avatarDstPoint := image.Pt(offsetX, offsetY)
	avatarRectangle := image.Rectangle{avatarDstPoint, avatarDstPoint.Add(rszAvatarImage.Bounds().Size())}
	draw.Draw(rgbaImage, avatarRectangle, rszAvatarImage,
		rszAvatarImage.Bounds().Min, draw.Over)

	// Draw the name date and message content.
	nameDrawer := newDrawer(rgbaImage, fontBold, nameColor, 18)
	nameX := offsetX + 50 + 10 // image width is 50px and offset is 10px
	nameY := offsetY + 18      // font size is 18px
	nameDrawer.Dot = fixed.P(nameX, nameY)
	nameDrawer.DrawString(*name)

	dateDrawer := newDrawer(rgbaImage, fontRegular, dateColor, 18)
	dateDrawer.Dot = fixed.Point26_6{
		// offset is 10px
		X: fixed.I(nameX+10) + nameDrawer.MeasureString(*name),
		// font size is 18px
		Y: fixed.I(offsetY + 18),
	}
	dateDrawer.DrawString(*date)

	contentDrawer := newDrawer(rgbaImage, fontBold, color.Black, 40)
	for i, line := range strings.Split(*content, "\\n") {
		// font size is 40px, spacing of lines is 3px
		contentDrawer.Dot = fixed.P(nameX, nameY+(40+3)*(i+1))
		contentDrawer.DrawString(line)
	}

	var outFile *os.File
	if *stdout {
		// If `stdout` flag is set then just output the image to the stdout.
		outFile = os.Stdout
	} else {
		if len(*outputName) == 0 {
			errorReturn("output file name MUST NOT be empty")
		}
		outFile, err = os.Create(*outputName + ".jpeg")
		if err != nil {
			errorReturn("create file failed: %v", err)
		}
	}
	defer outFile.Close()

	bufferWriter := bufio.NewWriter(outFile)
	err = jpeg.Encode(bufferWriter, rgbaImage, nil)
	if err != nil {
		errorReturn("encode image failed: %v", err)
	}
	err = bufferWriter.Flush()
	if err != nil {
		errorReturn("flush buffer to disk failed: %v", err)
	}

	if *isOpen && !*stdout {
		err = open.Run("out.jpeg")
		if err != nil {
			errorReturn("open image failed: %v", err)
		}
	}
}

// Make an image reader from a local image file or fetch from an remote url.
func makeImageReader() io.ReadCloser {
	if len(*imageSource) == 0 {
		errorReturn("avatar image source MUST NOT be empty")
	}

	if url, err := url.ParseRequestURI(*imageSource); err == nil {
		if len(url.Scheme) > 0 {
			// `imageSource` can be parsed and has a scheme, so it is a url.
			resp, err := http.Get(*imageSource)
			if err != nil {
				errorReturn("failed to fetch image: %v", err)
			}
			return resp.Body
		}
	}

	// otherwise read from local file system
	imageSourcePath, err := filepath.Abs(*imageSource)
	if err != nil {
		errorReturn("get image path failed: %v", err)
	}
	imageReader, err := os.Open(imageSourcePath)
	if err != nil {
		errorReturn("open image failed: %v", err)
	}
	return imageReader
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
			Size:    pixelsToPoints(fontSize),
			DPI:     *dpi,
			Hinting: font.HintingFull,
		}),
	}
}

// Convert pixels to points (font size uint)
// - `pt` / 72 * `DPI` = `px`
func pixelsToPoints(px float64) float64 {
	return px * 72 / *dpi
}

func errorReturn(f string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "\033[1;31m[error]\033[0m " + f, v...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}
