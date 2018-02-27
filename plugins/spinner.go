package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
)

func main() {
	outGif := &gif.GIF{}

	for i := 0; i < 4; i++ {
		dest := spinner(4)

		outGif.Image = append(outGif.Image, dest)
		outGif.Delay = append(outGif.Delay, 0)
	}
	// save to out.gif
	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)

	defer f.Close()
	gif.EncodeAll(f, outGif)

}

func spinner(seg int) *image.Paletted {
	var width float64 = 200
	var height float64 = 200
	bounds := image.Rect(0, 0, int(width), int(height))

	dest := image.NewRGBA(bounds)
	canvas := draw2dimg.NewGraphicContext(dest)

	for i := 0; i < seg; i++ {
		rand.Seed(time.Now().UnixNano() * int64(i+1*5))
		var startAng float64

		if i != 0 {
			startAng = 360 / float64(seg) * float64(i) * (math.Pi / 180.0)
		}

		var endAng = 360 / float64(seg) * (math.Pi / 180.0)

		fmt.Println("Start: " + strconv.FormatFloat(startAng, 'f', -1, 64) + " End: " + strconv.FormatFloat(endAng, 'f', -1, 64))
		var r = uint8(rand.Intn(255))
		var g = uint8(rand.Intn(255))
		var b = uint8(rand.Intn(255))

		fmt.Println("R:" + strconv.Itoa(int(r)) + " G:" + strconv.Itoa(int(g)) + " B:" + strconv.Itoa(int(b)))

		canvas.SetFillColor(color.RGBA{r, g, b, 255})
		canvas.SetStrokeColor(color.RGBA{r, g, b, 255})
		canvas.SetLineWidth(1)

		canvas.BeginPath()
		canvas.ArcTo(width/2, height/2, width/2, height/2, startAng, endAng)
		canvas.LineTo(width/2, height/2)
		canvas.Close()
		canvas.FillStroke()
	}

	palettedImage := image.NewPaletted(bounds, palette.Plan9)
	draw.Draw(palettedImage, palettedImage.Rect, dest, bounds.Min, draw.Over)

	return palettedImage
}
