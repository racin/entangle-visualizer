package main

import (
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/racin/entangle/entangler"
)

func addParityBetweenDatablock(img *ebiten.Image, dataLeft, dataRight int, fill color.Color, width, columnOffset int) {
	var dataLeftRow, dataRightRow float64
	if dataLeftRow = float64(dataLeft % entangler.HorizontalStrands); dataLeftRow == 0 {
		dataLeftRow = entangler.HorizontalStrands
	}
	if dataRightRow = float64(dataRight % entangler.HorizontalStrands); dataRightRow == 0 {
		dataRightRow = entangler.HorizontalStrands
	}
	var dataLeftColumn float64 = float64(int((dataLeft-1)/entangler.HorizontalStrands) + columnOffset)
	var dataRightColumn float64 = float64(int((dataRight-1)/entangler.HorizontalStrands) + columnOffset)
	var dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos float64

	/*
		row = i % entangler.HorizontalStrands
		column = int(i / entangler.HorizontalStrands)

		x := float64(xOffset + (xSpace * column))
		y := float64(yOffset + (ySpace * (row)))
	*/
	if dataLeftRow == dataRightRow { // Horizontal
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + dataRadius
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 1
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1))
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow+1 == dataRightRow { // Right
		dataLeftXpos = xOffset + dataLeftColumn*xSpace
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = dataLeftYpos + ySpace
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow-1 == dataRightRow { // Left
		dataLeftXpos = xOffset + dataLeftColumn*xSpace
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = dataLeftYpos - ySpace
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow > dataRightRow { // Need to draw two lines. Wrap left.
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + dataRadius
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = yOffset

		addEdge(img, dataLeftXpos-dataRadius, dataLeftYpos, dataLeftXpos+10, dataLeftYpos+dataRadius+10, fill, width)
		addEdge(img, dataRightXpos-(2*dataRadius)-10, dataRightYpos-dataRadius-10, dataRightXpos-dataRadius, dataRightYpos, fill, width)
	} else if dataLeftRow < dataRightRow { // Need to draw two lines. Wrap right.
		dataLeftXpos = xOffset + dataLeftColumn*xSpace
		dataLeftYpos = yOffset

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = yOffset + (ySpace * 4)

		addEdge(img, dataLeftXpos, dataLeftYpos, dataLeftXpos+dataRadius+10, dataLeftYpos-dataRadius-10, fill, width)
		addEdge(img, dataRightXpos-dataRadius-10, dataRightYpos+dataRadius+10, dataRightXpos, dataRightYpos, fill, width)
	}
}

func addDataBlock(img *ebiten.Image, radius float64, edge, fill, textColor color.Color, index, columnOffset int) {
	var row, column int
	i := index - 1
	row = i % entangler.HorizontalStrands
	column = int(i/entangler.HorizontalStrands) - columnOffset
	// if column < 0 {
	// 	column += int(lattice.NumDataBlocks/entangler.HorizontalStrands) + 1
	// }

	x := float64(xOffset + (xSpace * column))
	y := float64(yOffset + (ySpace * (row)))

	addCircle(img, x, y, radius, edge, fill)
	indexStr := strconv.Itoa(index)

	offset := 8.0 // (dia - dataFontSize) / 2
	xoffset := ((1 + math.Floor(math.Log10(float64(index)))) * offset) - 2

	text.Draw(img, indexStr, dataFont, int(x-xoffset), int(y+offset), textColor)
}
