package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type circleKey struct {
	edge color.Color
	fill color.Color
}

func newCircleKey(edge, fill color.Color) circleKey {
	return circleKey{edge: edge, fill: fill}
}

func slowAddCircle(img *ebiten.Image, x, y, radius float64, edge, fill color.Color) {
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
func addCircle(img *ebiten.Image, x, y, radius float64, edge, fill color.Color) {
	key := newCircleKey(edge, fill)
	if _, ok := circles[key]; !ok {
		circles[key], _ = ebiten.NewImage(dataRadius*2, dataRadius*2, ebiten.FilterDefault)
		slowAddCircle(circles[key], dataRadius, dataRadius, dataRadius, key.edge, key.fill)
	}
	g := ebiten.GeoM{}
	g.Translate(x-dataRadius, y-dataRadius)
	img.DrawImage(circles[key], &ebiten.DrawImageOptions{
		GeoM: g,
	})
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
	xdiff := endX - startX
	ydiff := endY - startY
	m := (ydiff) / (xdiff)
	red := color.RGBA{0xff, 0x0, 0x0, 0xff}

	var edge color.Color
	if fill == red {
		edge = fill
	} else {
		edge = color.Black
	}

	// for i := startY; i < endY; i++ {
	// 	if fill == red {
	// 		if math.Mod(math.Abs(i), 6) < 3 {
	// 			continue
	// 		}
	// 	}
	// 	a := int((ydiff-(endY-i))*m) + int(startX)
	// 	ii := int(i)
	// 	for j := 0; j < width; j++ {
	// 		if j <= 1 || j+2 >= width {
	// 			img.Set(ii, j+a, edge)
	// 		} else {
	// 			img.Set(ii, j+a, fill)
	// 		}
	// 	}
	// }

	if startX == endX {
		// fmt.Printf("Trying to decrement.... xdiff: %v, ydiff: %v, m: %v", xdiff, ydiff, m)
		start := startY
		end := endY
		if end < start {
			start = endY
			end = startY
		}
		for i := start; i < end; i++ {
			if fill == red {
				if math.Mod(math.Abs(i), 6) < 3 {
					continue
				}
			}
			ii := int(i)
			xx := int(startX)
			for j := 0; j < width; j++ {
				if i-7 < start {
					if i-2 < start || j+2 >= width {
						img.Set(xx+j, ii, edge)
					} else {
						img.Set(xx+j, ii, fill)
					}

				} else if j <= 1 || j+2 >= width {
					img.Set(xx+j, ii, edge)
				} else {
					img.Set(xx+j, ii, fill)
				}
			}
		}
	}
	for i := startX; i < endX; i++ {
		if fill == red {
			if math.Mod(math.Abs(i), 6) < 3 {
				continue
			}
		}
		a := int((xdiff-(endX-i))*m) + int(startY)
		ii := int(i)
		for j := 0; j < width; j++ {
			if j <= 1 || j+2 >= width {
				img.Set(ii, j+a, edge)
			} else {
				img.Set(ii, j+a, fill)
			}
		}
	}
}
