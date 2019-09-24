package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
	"log"
)

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
	screen.DrawImage(yy, nil)
	return nil
}

func main() {
	ebiten.SetMaxTPS(3)
	if err := ebiten.Run(update, 640, 480, 1, "Fill"); err != nil {
		log.Fatal(err)
	}
}
