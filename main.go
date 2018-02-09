package main

import (
	"image"
	"image/draw"
	"image/png"
	"os"
	"log"
	"github.com/skratchdot/open-golang/open"
	"bufio"
)

const (
	imageWidth = 500
	imageHeight = 200

	offsetPointX = 20
	offsetPointY = 10
)

func main() {
	// New a rgba image with 500px width and 200px height.
	rgbaImage := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Draw a background on the rgba image.
	backgroundImage := image.White
	draw.Draw(rgbaImage, rgbaImage.Bounds(), backgroundImage, image.ZP, draw.Src)

	// Draw the avatar (test image).
	avatarImageReader, err := os.Open("testdata/rsz_telegram.png")
	if err != nil {
		log.Fatal("open avatar image failed", err)
	}
	defer avatarImageReader.Close()

	avatarImage, _, err := image.Decode(avatarImageReader)
	if err != nil {
		log.Fatal("decode avatar image failed", err)
	}

	startPoint := image.Pt(offsetPointX, offsetPointY)
	draw.Draw(rgbaImage, rgbaImage.Bounds(), avatarImage, rgbaImage.Bounds().Min.Sub(startPoint), draw.Src)

	// Save image to disk.
	outFile, err := os.Create("out.png")
	if err != nil {
		log.Fatal("creat file failed", err)
	}
	defer outFile.Close()

	bufferWriter := bufio.NewWriter(outFile)
	err = png.Encode(bufferWriter, rgbaImage)
	if err != nil {
		log.Fatal("encode image failed", err)
	}
	err = bufferWriter.Flush()
	if err != nil {
		log.Fatal("flush buffer to disk failed", err)
	}

	// Open image.
	err = open.Run("out.png")
	if err != nil {
		log.Fatal("open image failed", err)
	}
}
