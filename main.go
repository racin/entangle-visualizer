package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/racin/entangle-visualizer/resources/fonts"
	"github.com/racin/entangle/entangler"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"math"
)

var (
	dataFont font.Face
	lattice  *entangler.Lattice
)

const (
	windowTitle  = "Entangle Visualizer"
	windowYSize  = 430
	windowXSize  = 1840
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

	lattice = entangler.NewLattice(3, 5, 5, "lattice.json", nil)
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

	numDatablocks := lattice.NumDataBlocks

	for i := 0; i < numDatablocks; i++ {
		row := math.Floor(float64(i) / float64(dataPrRow))
		num := (((i * entangler.HorizontalStrands) + int(row)) % numDatablocks) + 1
		addDataBlock(screen, dataRadius, color.Black, nil, color.White, num)
	}

	addDataBlock(screen, dataRadius, color.Black, color.RGBA{0, 0xff, 0, 0}, color.White, 12)
	addDataBlock(screen, dataRadius, color.Black, color.RGBA{0, 0xff, 0, 0}, color.White, 27)
	addDataBlock(screen, dataRadius, color.Black, color.RGBA{0, 0xff, 0, 0}, color.White, 1)
	addDataBlock(screen, dataRadius, color.Black, color.RGBA{0, 0xff, 0, 0}, color.White, 112)
	addDataBlock(screen, dataRadius, color.Black, color.RGBA{0, 0xff, 0, 0}, color.White, 18)
	addDataBlock(screen, dataRadius, color.Black, color.RGBA{0, 0xff, 0, 0}, color.White, 24)
	addDataBlock(screen, dataRadius, color.Black, color.RGBA{0, 0xff, 0, 0}, color.White, 30)

	for i := 0; i < len(lattice.Blocks); i++ {
		block := lattice.Blocks[i]
		if !block.IsParity {
			continue
		}
		var leftPos, rightPos int

		if len(block.Left) == 0 || block.Left[0].Position < 1 {
			rightPos = block.Right[0].Position + numDatablocks
			r, h, l := entangler.GetBackwardNeighbours(rightPos, entangler.S, entangler.P)
			switch block.Class {
			case entangler.Horizontal:
				leftPos = h
			case entangler.Right:
				leftPos = r
			case entangler.Left:
				leftPos = l
			}
		} else if len(block.Right) == 0 || block.Right[0].Position > numDatablocks+5 {
			continue
		} else {
			leftPos = block.Left[0].Position
			rightPos = block.Right[0].Position
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
		fmt.Printf("Print parity. Left: %d, Right: %d\n", leftPos, rightPos)
		addParityBetweenDatablock(screen, leftPos, rightPos, clr, 3)
	}

	return nil
}

func main() {
	ebiten.SetMaxTPS(3)
	if err := ebiten.Run(update, windowXSize, windowYSize, 1, windowTitle); err != nil {
		log.Fatal(err)
	}
}
