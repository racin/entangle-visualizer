package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/racin/entangle/entangler"
	"image/color"
	"math"
	"strconv"
)

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
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) + dataRadius + (ySpace * dataLeftColumn) - 2

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) + dataRadius + dataRadius + 3 + (ySpace * dataLeftColumn)
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

func addDataBlock(img *ebiten.Image, radius float64, edge, fill, textColor color.Color, index int) {
	var row, column int
	i := index - 1
	row = i % entangler.HorizontalStrands
	column = int(i / entangler.HorizontalStrands)
	// num := (((i * entangler.HorizontalStrands) + int(row)) % numDatablocks) + 1
	x := float64(xOffset + (xSpace * column))
	y := float64(yOffset + (ySpace * (row)))
	addCircle(img, x, y, radius, edge, fill)
	indexStr := strconv.Itoa(index)
	// x - 6, y + 6
	//dia := 2 * radius
	offset := 8.0 // (dia - dataFontSize) / 2
	xoffset := ((1 + math.Floor(math.Log10(float64(index)))) * offset) - 2
	//fmt.Printf("DIAMETER: %f, x: %f, xoffset: %f\n", dia, index, xoffset)
	text.Draw(img, indexStr, dataFont, int(x-xoffset), int(y+offset), textColor)
}
