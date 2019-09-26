package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/racin/entangle-visualizer/resources/fonts"
	"github.com/racin/entangle/entangler"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"log"
	"math"
	"strconv"
)

var (
	dataFont     font.Face
	dataFontSize float64
	dataFontDPI  float64
	dataPrRow    int
)

func init() {
	tt, err := truetype.Parse(fonts.OpenSans_Regular_tff)
	if err != nil {
		log.Fatal(err)
	}
	dataFontSize = 24
	dataFontDPI = 72
	dataPrRow = 20
	dataFont = truetype.NewFace(tt, &truetype.Options{
		Size:    dataFontSize,
		DPI:     dataFontDPI,
		Hinting: font.HintingFull,
	})
}
func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.Fill(color.RGBA{0xff, 0, 0, 0xff})
	fmt.Printf("Bounds. Min X: %d, Max X: %d, Min Y: %d, Max Y: %d\n",
		screen.Bounds().Min.X, screen.Bounds().Max.X,
		screen.Bounds().Min.Y, screen.Bounds().Max.Y)
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())
	text.Draw(screen, msg, basicfont.Face7x13, 200, 300, color.White)

	for i := 1; i <= 3; i++ {
		row := math.Floor(float64(i) / float64(dataPrRow))
		ix := (60 * (i % dataPrRow))
		addDataBlock(screen, float64(40+(ix)), float64(50+(60*(row))), 25, color.Black, nil, color.White, i)
	}

	addParity(screen, 10, 90, 120, 800, color.Black, 2)

	return nil
}

func addParity(img *ebiten.Image, startX, startY, endX, endY float64, fill color.Color, width int) {
	m := (endY - startY) / (endX - startX)
	for i := startX; i < endX; i++ {
		a := int(i * m)
		ii := int(i)
		for j := 0; j < width; j++ {
			img.Set(ii, int(startY)+j+a, fill)
		}
	}
}

func addDataBlock(img *ebiten.Image, x, y, radius float64, edge, fill, textColor color.Color, index int) {
	addCircle(img, x, y, radius, edge, fill)
	i := strconv.Itoa(index)
	// x - 6, y + 6
	dia := 2 * radius
	offset := 8.0 // (dia - dataFontSize) / 2
	xoffset := (1 + math.Floor(math.Log10(float64(index)))) * offset
	fmt.Printf("DIAMETER: %f, x: %f, xoffset: %f\n", dia, index, xoffset)
	text.Draw(img, i, dataFont, int(x-xoffset), int(y+offset), textColor)
}
func addCircle(img *ebiten.Image, x, y, radius float64, edge, fill color.Color) {
	var r2 float64 = radius * radius
	var i, j float64
	for i = -radius + 1; i < radius; i++ {
		for j = -radius + 1; j < radius; j++ {
			point := math.Pow(i, 2) + math.Pow(j, 2)
			if point <= r2 {
				if math.Abs(point-r2) < 2*radius {
					img.Set(int(i+x), int(j+y), edge)
				} else if fill != nil {
					img.Set(int(i+x), int(j+y), fill)
				}
			}
		}
	}

}

func addSquare(img *ebiten.Image, x, y, length, width float64, fill bool) {
	var i, j float64
	for i = 0; i < length; i++ {
		for j = 0; j < width; j++ {
			if fill || i == 0 || j == 0 || i == length-1 || j == width-1 {
				img.Set((int)(i+x), (int)(j+y), color.RGBA{0, 0, 0, 0xff})
			}
		}
	}

}

func main() {
	ebiten.SetMaxTPS(3)
	if err := ebiten.Run(update, 1400, 900, 1, "Fill"); err != nil {
		log.Fatal(err)
	}
}
