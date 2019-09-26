package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/racin/entangle-visualizer/resources/fonts"
	"github.com/racin/entangle/entangler"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"math"
	"strconv"
)

var (
	dataFont font.Face
	lattice  *entangler.Lattice
)

const (
	dataFontSize = 24.0
	dataFontDPI  = 72.0
	dataPrRow    = 23 // 23
	xSpace       = 80.0
	xOffset      = 40.0
	ySpace       = 80.0
	yOffset      = 50.0
	dataRadius   = 25.0
)

func init() {
	tt, err := truetype.Parse(fonts.OpenSans_Regular_tff)
	if err != nil {
		log.Fatal(err)
	}

	lattice = entangler.NewLattice(3, 5, 5, "../entangle/retrives.txt", nil)
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

	numDatablocks := 115

	for i := 0; i < numDatablocks; i++ {
		row := math.Floor(float64(i) / float64(dataPrRow))
		ix := (xSpace * (i % dataPrRow))
		num := (((i * entangler.HorizontalStrands) + int(row)) % numDatablocks) + 1
		addDataBlock(screen, float64(xOffset+(ix)), float64(yOffset+(ySpace*(row))), dataRadius, color.Black, nil, color.White, num)
	}

	for i := 0; i < len(lattice.Blocks); i++ {
		block := lattice.Blocks[i]
		if !block.IsParity {
			continue
		}
		if len(block.Left) == 0 || len(block.Right) == 0 || block.Left[0].Position < 1 || block.Right[0].Position > numDatablocks {
			continue
		}
		// if block.Right[0].Position > numDatablocks && (block.Right[0].Position == block.Left[0].Position+1 || block.Right[0].Position == block.Left[0].Position+9) {
		// 	continue
		// }
		var clr color.Color
		switch block.Class {
		case entangler.Horizontal:
			clr = color.RGBA{0, 0xff, 0, 0xff}
		case entangler.Right:
			clr = color.RGBA{0, 0, 0xff, 0xff}
		case entangler.Left:
			clr = color.Black
		}
		fmt.Printf("Print parity. Left: %d, Right: %d\n", block.Left[0].Position, block.Right[0].Position)
		addParityBetweenDatablock(screen, block.Left[0].Position, block.Right[0].Position, clr, 3)
	}

	//addParityBetweenDatablock(screen, 1, 6, color.Black, 2)
	//addParityBetweenDatablock(screen, 6, 11, color.Black, 2)
	//addParityBetweenDatablock(screen, 1, 7, color.Black, 2)
	//addParityBetweenDatablock(screen, 3, 7, color.RGBA{0, 0xff, 0, 0}, 2)

	//addParity(screen, xOffset, yOffset+ySpace, xOffset+xSpace, yOffset+ySpace, color.White, 2)

	return nil
}

func addParityBetweenDatablock(img *ebiten.Image, dataLeft, dataRight int, fill color.Color, width int) {
	var dataLeftRow, dataRightRow float64
	if dataLeftRow = float64(dataLeft % entangler.HorizontalStrands); dataLeftRow == 0 {
		dataLeftRow = entangler.HorizontalStrands
	}
	if dataRightRow = float64(dataRight % entangler.HorizontalStrands); dataRightRow == 0 {
		dataRightRow = entangler.HorizontalStrands
	}
	var dataLeftColumn float64 = float64(int((dataLeft - 1) / entangler.HorizontalStrands))
	var dataRightColumn float64 = float64(int((dataRight - 1) / entangler.HorizontalStrands))
	var dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos float64

	if dataLeftRow == dataRightRow {
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + dataRadius
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 1
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1))
		addParity(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow+1 == dataRightRow {
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + 8
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) - dataRadius - (ySpace * dataLeftColumn)

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) - dataRadius - dataRadius - 5 - (ySpace * dataLeftColumn)
		addParity(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
		fmt.Printf("Printing Right. StartX: %f, StartY: %f, EndX: %f, EndY: %f, DataLeftRow: %f\n", dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, dataLeftRow)
	} else if dataLeftRow-1 == dataRightRow {
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + 8
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) + dataRadius + (ySpace * dataLeftColumn)

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) + dataRadius + dataRadius + 5 + (ySpace * dataLeftColumn)
		addParity(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
		fmt.Printf("Printing Left. StartX: %f, StartY: %f, EndX: %f, EndY: %f\n", dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos)
	} else if dataLeftRow > dataRightRow { // Need to draw two lines
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + 8
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) - dataRadius - (ySpace * dataLeftColumn)

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) + 5 - dataRadius - dataRadius - dataRadius - dataRadius - (ySpace * dataLeftColumn)

		addParity(img, dataLeftXpos, dataLeftYpos, dataLeftXpos+15, dataLeftYpos+15, fill, width)
		addParity(img, dataRightXpos-15, dataRightYpos-15, dataRightXpos, dataRightYpos, fill, width)
		fmt.Printf("Printing WRAPPER Right. StartX: %f, StartY: %f, EndX: %f, EndY: %f\n", dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos)
	} else if dataLeftRow < dataRightRow { // Need to draw two lines
		dataLeftXpos = xOffset + dataLeftColumn*xSpace
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) + 15 + (ySpace * dataLeftColumn)

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) + dataRadius + dataRadius + dataRadius + dataRadius + 20 + (ySpace * dataLeftColumn)

		addParity(img, dataLeftXpos, dataLeftYpos, dataLeftXpos+15, dataLeftYpos-15, fill, width)
		addParity(img, dataRightXpos, dataRightYpos+15, dataRightXpos+15, dataRightYpos, fill, width)
		fmt.Printf("Printing WRAPPER Left. StartX: %f, StartY: %f, EndX: %f, EndY: %f\n", dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos)
	}
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
	//dia := 2 * radius
	offset := 8.0 // (dia - dataFontSize) / 2
	xoffset := ((1 + math.Floor(math.Log10(float64(index)))) * offset) - 2
	//fmt.Printf("DIAMETER: %f, x: %f, xoffset: %f\n", dia, index, xoffset)
	text.Draw(img, i, dataFont, int(x-xoffset), int(y+offset), textColor)
}
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

func main() {
	ebiten.SetMaxTPS(3)
	if err := ebiten.Run(update, 1850, 900, 1, "Fill"); err != nil {
		log.Fatal(err)
	}
}
