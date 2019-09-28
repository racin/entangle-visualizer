package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math"
)

func addCircle(img *ebiten.Image, x, y, radius float64, edge, fill color.Color) {
	var r2 float64 = radius * radius
	var i, j float64
	for i = -radius + 1; i < radius; i++ {
		for j = -radius + 1; j < radius; j++ {
			point := math.Pow(i, 2) + math.Pow(j, 2)
			if point <= r2 {
				if math.Abs(point-r2) < 4*radius {
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

func addEdge(img *ebiten.Image, startX, startY, endX, endY float64, fill color.Color, width int) {
	m := (endY - startY) / (endX - startX)
	for i := startX; i < endX; i++ {
		a := int(i * m)
		ii := int(i)
		for j := 0; j < width; j++ {
			img.Set(ii, int(startY)+j+a, fill)
		}
	}
}
