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

var width float64 = 256
var height float64 = 256
var bounds = image.Rect(0, 0, int(width), int(height))
var fileName = "/tmp/out.gif"

func send(d *discordgo.Session, m *discordgo.MessageCreate, fileName string, timen time.Time) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		d.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	ms := &discordgo.MessageSend{
		Content: time.Since(timen).String(),
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

func Spinner(d *discordgo.Session, m *discordgo.MessageCreate, num int) {
	timen := time.Now()

	outGif := &gif.GIF{}
	//d.ChannelMessageSend(m.ChannelID, "start")
	wheel := wheel(num)
	//d.ChannelMessageSend(m.ChannelID, "made wheel")

	for i := 0; i < 30; i++ {
		//go d.ChannelMessageSend(m.ChannelID, "loop number "+strconv.Itoa(i))
		//timel := time.Now()
		round := rotateAround(wheel, float64(12*(i+1)), int(width/2), int(height/2))
		//times := time.Since(timel)
		//go d.ChannelMessageSend(m.ChannelID, times.String()+"rotated")

		//timel = time.Now()
		//dest := image.NewPaletted(bounds, palette.Plan9)
		//draw.Draw(dest, dest.Rect, round, bounds.Min, draw.Over)
		//times = time.Since(timel)
		//go d.ChannelMessageSend(m.ChannelID, times.String()+"encoded")

		//timel = time.Now()
		outGif.Image = append(outGif.Image, round)
		outGif.Delay = append(outGif.Delay, 0)
		//times = time.Since(timel)
		//go d.ChannelMessageSend(m.ChannelID, times.String()+"appended")
	}

	// save to out.gif
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)

	defer f.Close()
	gif.EncodeAll(f, outGif)

	send(d, m, fileName, timen)
}

func rotateAround(src image.Image, angle float64, x, y int) *image.Paletted {
	a := math.Pi * angle / 180
	xf, yf := float64(x), float64(y)
	sin := math.Sin(a)
	cos := math.Cos(a)
	matrix := f64.Aff3{
		cos, -sin, xf - xf*cos + yf*sin,
		sin, cos, yf - xf*sin - yf*cos,
	}
	dst := image.NewPaletted(src.Bounds(), palette.Plan9)
	draw.BiLinear.Transform(dst, matrix, src, src.Bounds(), draw.Src, nil)
	return dst
}

func wheel(seg int) *image.Paletted {
	//dest := image.NewRGBA(bounds)
	dest := image.NewPaletted(bounds, palette.Plan9)
	canvas := draw2dimg.NewGraphicContext(dest)

	for i := 0; i < seg; i++ {
		rand.Seed(time.Now().UnixNano() * int64(i+1*5))
		var (
			r        = uint8(rand.Intn(255))
			g        = uint8(rand.Intn(255))
			b        = uint8(rand.Intn(255))
			startAng float64
			endAng   = 360 / float64(seg) * (math.Pi / 180.0)
		)

		if i != 0 {
			startAng = 360 / float64(seg) * float64(i) * (math.Pi / 180.0)
		}

		canvas.SetFillColor(color.RGBA{r, g, b, 255})
		canvas.SetStrokeColor(color.RGBA{r, g, b, 255})
		canvas.SetLineWidth(1)

		canvas.BeginPath()
		canvas.ArcTo(width/2, height/2, width/2, height/2, startAng, endAng)
		canvas.LineTo(width/2, height/2)
		canvas.Close()
		canvas.FillStroke()
	}
	return dest
}
