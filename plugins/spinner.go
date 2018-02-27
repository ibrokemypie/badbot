package plugins

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/llgcode/draw2d/draw2dimg"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

var width float64 = 200
var height float64 = 200
var bounds = image.Rect(0, 0, int(width), int(height))
var fileName = "/tmp/out.gif"

func Spinner(d *discordgo.Session, m *discordgo.MessageCreate, num int) {
	timen := time.Now()

	outGif := &gif.GIF{}
	dest := wheel(num)

	for i := 0; i < 6; i++ {
		dest = rotate(dest)

		outGif.Image = append(outGif.Image, dest)
		outGif.Delay = append(outGif.Delay, 0)
	}

	// save to out.gif
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)

	defer f.Close()
	gif.EncodeAll(f, outGif)
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		d.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	time := time.Since(timen)

	ms := &discordgo.MessageSend{
		Content: time.String(),
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + fileName,
			},
		},
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   fileName,
				Reader: f,
			},
		},
	}
	d.ChannelMessageSendComplex(m.ChannelID, ms)

}

func rotateAround(src image.Image, angle float64, x, y int) image.Image {
	a := math.Pi * angle / 180
	xf, yf := float64(x), float64(y)
	sin := math.Sin(a)
	cos := math.Cos(a)
	matrix := f64.Aff3{
		cos, -sin, xf - xf*cos + yf*sin,
		sin, cos, yf - xf*sin - yf*cos,
	}
	dst := image.NewRGBA(src.Bounds())
	draw.BiLinear.Transform(dst, matrix, src, src.Bounds(), draw.Src, nil)
	return dst
}

func rotate(in *image.Paletted) *image.Paletted {

	dest := rotateAround(in, 60, int(width/2), int(height/2))
	palettedImage := image.NewPaletted(bounds, palette.Plan9)
	draw.Draw(palettedImage, palettedImage.Rect, dest, bounds.Min, draw.Over)

	return palettedImage
}

func wheel(seg int) *image.Paletted {
	dest := image.NewRGBA(bounds)
	canvas := draw2dimg.NewGraphicContext(dest)

	for i := 0; i < seg; i++ {
		rand.Seed(time.Now().UnixNano() * int64(i+1*5))
		var startAng float64

		if i != 0 {
			startAng = 360 / float64(seg) * float64(i) * (math.Pi / 180.0)
		}

		var endAng = 360 / float64(seg) * (math.Pi / 180.0)

		//fmt.Println("Start: " + strconv.FormatFloat(startAng, 'f', -1, 64) + " End: " + strconv.FormatFloat(endAng, 'f', -1, 64))
		var r = uint8(rand.Intn(255))
		var g = uint8(rand.Intn(255))
		var b = uint8(rand.Intn(255))

		//fmt.Println("R:" + strconv.Itoa(int(r)) + " G:" + strconv.Itoa(int(g)) + " B:" + strconv.Itoa(int(b)))

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
