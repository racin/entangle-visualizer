package main

import (
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
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow+1 == dataRightRow {
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + 8
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) - dataRadius - (ySpace * dataLeftColumn)

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) - dataRadius - dataRadius - 5 - (ySpace * dataLeftColumn)
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow-1 == dataRightRow {
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + 8
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) + dataRadius + (ySpace * dataLeftColumn) - 2

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) + dataRadius + dataRadius + 3 + (ySpace * dataLeftColumn)
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow > dataRightRow { // Need to draw two lines
		dataLeftXpos = xOffset + dataLeftColumn*xSpace + 8
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) - dataRadius - (ySpace * dataLeftColumn)

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) + 5 - dataRadius - dataRadius - dataRadius - dataRadius - (ySpace * dataLeftColumn)

		addEdge(img, dataLeftXpos, dataLeftYpos, dataLeftXpos+15, dataLeftYpos+15, fill, width)
		addEdge(img, dataRightXpos-15, dataRightYpos-15, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow < dataRightRow { // Need to draw two lines
		dataLeftXpos = xOffset + dataLeftColumn*xSpace
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1)) + 15 + (ySpace * dataLeftColumn)

		dataRightXpos = xOffset + dataRightColumn*xSpace - dataRadius + 3
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1)) + dataRadius + dataRadius + dataRadius + dataRadius + 20 + (ySpace * dataLeftColumn)

		addEdge(img, dataLeftXpos, dataLeftYpos, dataLeftXpos+15, dataLeftYpos-15, fill, width)
		addEdge(img, dataRightXpos, dataRightYpos+15, dataRightXpos+15, dataRightYpos, fill, width)
	}
}

func addDataBlock(img *ebiten.Image, radius float64, edge, fill, textColor color.Color, index int) {
	var row, column int
	i := index - 1
	row = i % entangler.HorizontalStrands
	column = int(i / entangler.HorizontalStrands)

	x := float64(xOffset + (xSpace * column))
	y := float64(yOffset + (ySpace * (row)))
	addCircle(img, x, y, radius, edge, fill)
	indexStr := strconv.Itoa(index)

	offset := 8.0 // (dia - dataFontSize) / 2
	xoffset := ((1 + math.Floor(math.Log10(float64(index)))) * offset) - 2
	text.Draw(img, indexStr, dataFont, int(x-xoffset), int(y+offset), textColor)
}
