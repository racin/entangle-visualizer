package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/racin/entangle-visualizer/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
)

var (
	mplusNormalFont font.Face
)

func init() {
	basicfont.Face7x13.Width = 40
	basicfont.Face7x13.Height = 50
	tt, err := truetype.Parse(fonts.OpenSans_Regular_tff)
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
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
	y := screen.SubImage(image.Rect(100, 100, 200, 300))
	yy, _ := ebiten.NewImageFromImage(y, ebiten.FilterDefault)
	yy.Fill(color.RGBA{0, 0, 0xff, 0xff})
	addCircle(screen, 150, 150, 40, color.RGBA{0, 0, 0, 0xff}, nil)
	addCircle(screen, 200, 100, 10, color.RGBA{0, 0, 0, 0xff}, color.RGBA{0, 0, 0xff, 0xff})
	addCircle(screen, 300, 100, 20, color.RGBA{0, 0, 0, 0xff}, color.RGBA{0, 0, 0xff, 0xff})
	addSquare(screen, 400, 300, 40, 40, false)
	addSquare(screen, 400, 500, 80, 40, true)
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS())

	text.Draw(screen, msg, basicfont.Face7x13, 200, 300, color.White)
	for i := 5; i < 10; i++ {
		addDataBlock(screen, float64(50+(90*i)), 150, 30, color.Black, nil, color.White, i)
	}

	screen.DrawImage(yy, nil)

	return nil
}

func addDataBlock(img *ebiten.Image, x, y, radius float64, edge, fill, textColor color.Color, index int) {
	addCircle(img, x, y, radius, edge, fill)
	i := strconv.Itoa(index)
	text.Draw(img, i, mplusNormalFont, int(x), int(y), textColor)
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
	if err := ebiten.Run(update, 1024, 768, 1, "Fill"); err != nil {
		log.Fatal(err)
	}
}
