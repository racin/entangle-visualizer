package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/racin/snarl/entangler"
)

var colors []color.Color = []color.Color{color.RGBA{0xff, 0x0, 0x0, 0xff}, color.RGBA{0x0, 0xff, 0, 0xff}, color.RGBA{0x33, 0x99, 0xff, 0xff},
	/* Green */ color.RGBA{0x0A, 0x66, 0x00, 0xFF}, color.RGBA{0x16, 0xE0, 0x00, 0xFF}, color.RGBA{0x0E, 0x8F, 0x00, 0xFF}, color.RGBA{0x12, 0xB8, 0x00, 0xFF}, color.RGBA{0x06, 0x3D, 0x00, 0xFF},
	/* Blue */ color.RGBA{0x00, 0x47, 0x8F, 0xFF}, color.RGBA{0x0A, 0x85, 0xFF, 0xFF}, color.RGBA{0x00, 0x5C, 0xB8, 0xFF}, color.RGBA{0x00, 0x70, 0xE0, 0xFF}, color.RGBA{0x00, 0x33, 0x66, 0xFF},
	/* Red */ color.RGBA{0xA3, 0x00, 0x00, 0xFF}, color.RGBA{0xFF, 0x47, 0x47, 0xFF}, color.RGBA{0xCC, 0x00, 0x00, 0xFF}, color.RGBA{0xFF, 0x1F, 0x1F, 0xFF}, color.RGBA{0x7A, 0x00, 0x00, 0xFF}}

func getColorForSpecialParities(oldClr color.Color, class entangler.StrandClass, rightRow int) color.Color {
	if oldClr == colors[0] {
		// Red failure line
		return colors[12+rightRow]

	} else if oldClr == colors[1] {
		// Green download line
		return colors[2+rightRow]

	} else {
		// Blue repaired line
		return colors[7+rightRow]
	}
}

// func isFillRed(fill color.Color) bool {
// 	reds := []int{0, 12, 13, 14, 15, 17}
// 	for i := 0; i < len(colors); i++ {
// 		if fill == colors[i] {
// 			for j := 0; j < len(reds); j++ {
// 				if i == reds[j] {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

func addParityBetweenDatablock(img *ebiten.Image, dataLeft, dataRight int, fill color.Color, width, columnOffset, horizontalStands int, class entangler.StrandClass) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("Recovered.... %v\n", rec)
		}
	}()

	var dataLeftRow, dataRightRow float64
	if dataLeftRow = float64(dataLeft % horizontalStands); dataLeftRow == 0 {
		dataLeftRow = float64(horizontalStands)
	}
	if dataRightRow = float64(dataRight % horizontalStands); dataRightRow == 0 {
		dataRightRow = float64(horizontalStands)
	}
	var dataLeftColumn float64 = float64(int((dataLeft-1)/horizontalStands) + columnOffset)
	var dataRightColumn float64 = float64(int((dataRight-1)/horizontalStands) + columnOffset)
	var dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos float64
	screenWidthFloat := float64(windowXSize) / zoom
	dataLeftXpos = xOffset + dataLeftColumn*xSpace
	if dataLeftXpos < 0 {
		dataLeftXpos += displayXSize
	}
	if dataLeftXpos > screenWidthFloat {
		return
	}
	/*
		row = i % entangler.HorizontalStrands
		column = int(i / entangler.HorizontalStrands)

		x := float64(xOffset + (xSpace * column))
		y := float64(yOffset + (ySpace * (row)))
	*/
	if dataLeft > dataRight {
		// fmt.Printf("Want to draw between %v and %v, but this is a special case, RightCol: %v, Class: %v, screenWidthFloat: %v\n", dataLeft, dataRight, dataRightColumn, class, screenWidthFloat)
		// screenWidth, _ := img.Size()
		if class == entangler.Horizontal {
			// Just draw a long line

			dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))
			dataRightYpos = dataLeftYpos

			dataRightXpos = xOffset + dataRightColumn*xSpace //dataLeftXpos + xSpace //
			if dataRightXpos < dataLeftXpos {
				dataRightXpos += displayXSize
			}
			// fmt.Printf("dataLeftXpos: %v, dataRightOrig: %v, dataRightXpos: %v, colOffset: %v, xOffSet: %v, dataLeftCol: %v, dataRightCol: %v, WindowXSize: %v, WindowXSize/Zoom: %v, xSpace: %v\n", dataLeftXpos, xOffset+dataRightColumn*xSpace, dataRightXpos, columnOffset, xOffset, dataLeftColumn, dataRightColumn, windowXSize, screenWidthFloat, xSpace)
			if dataRightColumn < 0 && dataRightXpos > screenWidthFloat {
				dataRightXpos = screenWidthFloat
			}

			addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
		} else if class == entangler.Right {
			// Draw 5 lines. South-East, East, North, East, South-East
			fillclr := getColorForSpecialParities(fill, class, int(dataRightRow))

			// South-East
			startX := dataLeftXpos
			startY := yOffset + (ySpace * (dataLeftRow - 1))
			endX := startX + (xSpace / 2) - 8
			endY := startY + (ySpace / 2) - 8

			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// East
			dataRightXpos = xOffset + dataRightColumn*xSpace //dataLeftXpos + xSpace //
			if dataRightXpos < dataLeftXpos {
				dataRightXpos += displayXSize
			}
			startX = endX
			startY = endY
			endX = dataRightXpos - 2.5*dataRadius + (float64(width-1) * dataRightRow)
			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// North
			startX = endX
			startY = yOffset + (ySpace * (dataRightRow - 2)) + (ySpace / 2) // - float64(width)
			endX = startX
			if startY < endY {
				endY += float64(width)
			}

			// fmt.Printf("StartX: %v, StartY: %v, EndX: %v, EndY: %v, ScreenWidth: %v\n", startX, startY, endX, endY, screenWidthFloat)
			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// East
			startX = endX
			endY = startY
			endX = dataRightXpos - dataRadius + 4 // - (xSpace / 2)
			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// South-East
			startX = endX
			startY = endY
			endX = dataRightXpos
			endY = yOffset + (ySpace * (dataRightRow - 1))
			if dataRightColumn < 0 && endX > screenWidthFloat {
				return
			}
			addEdge(img, startX, startY, endX, endY, fillclr, width)
		} else if class == entangler.Left {
			// Draw 5 lines. North-East, East, South, East, North-East
			fillclr := getColorForSpecialParities(fill, class, int(dataRightRow))

			// North-East
			startX := dataLeftXpos
			startY := yOffset + (ySpace * (dataLeftRow - 1))
			endX := startX + (xSpace / 2) - 8
			endY := startY - (ySpace / 2) - 8

			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// East
			dataRightXpos = xOffset + dataRightColumn*xSpace //dataLeftXpos + xSpace //
			if dataRightXpos < dataLeftXpos {
				dataRightXpos += displayXSize
			}
			startX = endX
			startY = endY
			endX = dataRightXpos - 2.5*dataRadius + (float64(width-1) * dataRightRow)
			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// South
			startX = endX
			startY = endY
			endX = startX
			endY = yOffset + (ySpace * (dataRightRow - 1)) + (ySpace / 2)
			// fmt.Printf("StartX: %v, StartY: %v, EndX: %v, EndY: %v\n", startX, startY, endX, endY)
			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// East
			startX = endX
			startY = endY
			endX = dataRightXpos - dataRadius //(xSpace / 2)
			addEdge(img, startX, startY, endX, endY, fillclr, width)
			// North-East
			startX = endX
			startY = endY
			endX = dataRightXpos
			endY = yOffset + (ySpace * (dataRightRow - 1))
			if dataRightColumn < 0 && endX > screenWidthFloat {
				return
			}
			addEdge(img, startX, startY, endX, endY, fillclr, width)
		}
	} else if dataLeftRow == dataRightRow { // Horizontal

		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = dataLeftXpos + xSpace // xOffset + dataRightColumn*xSpace + 1
		dataRightYpos = yOffset + (ySpace * (dataRightRow - 1))
		if dataRightColumn < 0 && dataRightXpos > screenWidthFloat {
			dataRightXpos = screenWidthFloat
		}
		// fmt.Printf("dataLeftXpos: %v, dataRightXpos: %v, LeftCol: %v, screenWidthFloat: %v\n", dataLeftXpos, dataRightXpos, dataLeftColumn, screenWidthFloat)
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow+1 == dataRightRow { // Right
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = dataLeftYpos + ySpace
		if dataRightColumn < 0 && dataRightXpos > screenWidthFloat {
			dataRightXpos = screenWidthFloat
		}
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow-1 == dataRightRow { // Left
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = dataLeftYpos - ySpace
		if dataRightColumn < 0 && dataRightXpos > screenWidthFloat {
			dataRightXpos = screenWidthFloat
		}
		addEdge(img, dataLeftXpos, dataLeftYpos, dataRightXpos, dataRightYpos, fill, width)
	} else if dataLeftRow > dataRightRow { // Need to draw two lines. Wrap Right.
		dataLeftYpos = yOffset + (ySpace * (dataLeftRow - 1))

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = yOffset

		addEdge(img, dataLeftXpos-dataRadius, dataLeftYpos, dataLeftXpos+10, dataLeftYpos+dataRadius+10, fill, width)
		if dataRightColumn < 0 && dataRightXpos > screenWidthFloat {
			return
		}
		addEdge(img, dataRightXpos-(2*dataRadius)-10, dataRightYpos-dataRadius-10, dataRightXpos-dataRadius, dataRightYpos, fill, width)
	} else if dataLeftRow < dataRightRow { // Need to draw two lines. Wrap Left.
		dataLeftYpos = yOffset

		dataRightXpos = dataLeftXpos + xSpace
		dataRightYpos = yOffset + (ySpace * 4)

		addEdge(img, dataLeftXpos, dataLeftYpos, dataLeftXpos+dataRadius+10, dataLeftYpos-dataRadius-10, fill, width)
		if dataRightColumn < 0 && dataRightXpos > screenWidthFloat {
			return
		}
		addEdge(img, dataRightXpos-dataRadius-10, dataRightYpos+dataRadius+10, dataRightXpos, dataRightYpos, fill, width)
	}
}

func addDataBlock(img *ebiten.Image, radius float64, edge, fill, textColor color.Color, index, columnOffset, horizontalStands int) {
	var row, column int
	i := index - 1
	row = i % horizontalStands
	column = int(i/horizontalStands) - columnOffset
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
