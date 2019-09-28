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
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	dataFont          font.Face
	lattice           *entangler.Lattice
	didDrawDataBlocks bool
	circles           map[circleKey]*ebiten.Image
	keyPresses        int
	keyLock           sync.Mutex
	keyLockBool       bool
	logParser         LogParser
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
	didDrawDataBlocks = false
	tt, err := truetype.Parse(fonts.OpenSans_Regular_tff)
	if err != nil {
		log.Fatal(err)
	}
	latticePath := "lattice.json"
	logPath := "output.txt"
	if _, err := os.Stat(latticePath); os.IsNotExist(err) {
		latticePath = "resources/lattice.json"
		if _, err := os.Stat(latticePath); os.IsNotExist(err) {
			fmt.Println("Did not find lattice.json in working directory or resources/")
			os.Exit(1)
		}
	}
	lattice = entangler.NewLattice(3, 5, 5, latticePath, nil)

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logPath = "resources/output.txt"
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			fmt.Println("Did not find output.txt in working directory or resources/")
			os.Exit(1)
		}
	}
	logParser := NewLogParser(logPath)
	if err := logParser.ReadLog(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	dataFont = truetype.NewFace(tt, &truetype.Options{
		Size:    dataFontSize,
		DPI:     dataFontDPI,
		Hinting: font.HintingFull,
	})
	circles = make(map[circleKey]*ebiten.Image)
	keyLockBool = false
}

func keyPressed(key ebiten.Key, presses int) {
	keyLock.Lock()
	defer keyLock.Unlock()
	if presses < keyPresses {
		return
	}

	fmt.Println("KEY PRESSED!! " + key.String())
	i, _ := strconv.Atoi(key.String())
	if lattice.Blocks[i].IsUnavailable {
		lattice.Blocks[i].Data = make([]byte, 5)
	} else {
		lattice.Blocks[i].IsUnavailable = true
	}

	time.Sleep(300 * time.Millisecond)
	keyPresses++
}

func update(screen *ebiten.Image) error {
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			go keyPressed(k, keyPresses)
			break
		}
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if !didDrawDataBlocks {
		didDrawDataBlocks = true
	}
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	numDatablocks := len(lattice.Blocks)
	for i := 0; i < numDatablocks; i++ {
		bl := lattice.Blocks[i]
		if bl.IsParity {
			continue
		}
		var clr color.Color
		if bl.HasData() {
			clr = color.RGBA{0x0, 0xff, 0x0, 0xff}
		} else if bl.IsUnavailable {
			clr = color.RGBA{0xff, 0x0, 0x0, 0xff}
		} else {
			clr = color.RGBA{0xc8, 0xc8, 0xc8, 0xff}
		}
		addDataBlock(screen, dataRadius, color.Black,
			clr, color.Black,
			lattice.Blocks[i].Position)
	}

	for i := 0; i < len(lattice.Blocks); i++ {
		block := lattice.Blocks[i]
		if !block.IsParity || !block.HasData() {
			continue
		}
		var leftPos, rightPos int

		if len(block.Left) == 0 || block.Left[0].Position < 1 {
			rightPos = block.Right[0].Position + lattice.NumDataBlocks
			r, h, l := entangler.GetBackwardNeighbours(rightPos, entangler.S, entangler.P)
			switch block.Class {
			case entangler.Horizontal:
				leftPos = h
			case entangler.Right:
				leftPos = r
			case entangler.Left:
				leftPos = l
			}
		} else if len(block.Right) == 0 || block.Right[0].Position > lattice.NumDataBlocks+5 {
			continue
		} else {
			leftPos = block.Left[0].Position
			rightPos = block.Right[0].Position
		}

		var clr color.Color = color.Black
		// switch block.Class {
		// case entangler.Horizontal:
		// 	clr = color.RGBA{0, 0xff, 0, 0xff}
		// case entangler.Right:
		// 	clr = color.RGBA{0, 0, 0xff, 0xff}
		// case entangler.Left:
		// 	clr = color.Black
		// }

		addParityBetweenDatablock(screen, leftPos, rightPos, clr, 3)
	}

	return nil
}

func main() {
	ebiten.SetMaxTPS(60)
	ebiten.SetRunnableInBackground(true)
	if err := ebiten.Run(update, windowXSize, windowYSize, 1, windowTitle); err != nil {
		log.Fatal(err)
	}
}
