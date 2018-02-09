package main

import (
	"image"
	_ "image/png"
	"image/draw"
	"os"
	"log"
	"github.com/skratchdot/open-golang/open"
)

const (
	imageWidth = 500
	imageHeight = 200

	offzetPointX = 20
	offzetPointY = 10
)

func main() {
	// New a rgba image with 500px width and 200px height.
	rgbaImage := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Draw a background on the rgba image.
	backgroundImage := image.White
	draw.Draw(rgbaImage, rgbaImage.Bounds(), backgroundImage, image.ZP, draw.Src)

	// Draw the avatar (test image)
	avatarImageReader, err := os.Open("testdata/rsz_telegram.png")
	if err != nil {
		log.Fatal("open avatar image failed", err)
	}
	defer avatarImageReader.Close()

	avatarImage, _, err := image.Decode(avatarImageReader)
	if err != nil {
		log.Fatal("decode avatar image failed", err)
	}

	draw.Draw(rgbaImage, rgbaImage.Bounds(), avatarImage, image.Point{offzetPointX, offzetPointY}, draw.Src)

}
