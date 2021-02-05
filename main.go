package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/racin/entangle-visualizer/resources/fonts"
	"github.com/racin/snarl/entangler"
	"golang.org/x/image/font"
)

var (
	dataFont     font.Face
	lattice      *entangler.Lattice
	circles      map[circleKey]*ebiten.Image
	keyPresses   int
	keyLock      sync.Mutex
	keyLockBool  bool
	logParser    *LogParser
	columnOffset int
	columninc    int     = 4
	zoom         float64 = 0.6
	zoominc      float64 = 0.2
	windowYSize  int     = 500
	windowXSize  int
	displayXSize float64
	logfile      string
)

const (
	windowTitle = "Entangle Visualizer"

	dataFontSize = 24.0
	dataFontDPI  = 72.0
	//dataPrRow    = 23 // 23
	xSpace     = 80.0
	xOffset    = 40.0
	ySpace     = 80.0
	yOffset    = 50.0
	dataRadius = 25.0
)

func init() {
	if len(os.Args) < 2 {
		fmt.Printf("Must supply path to log file.")
		os.Exit(1)
	}
	logfile = os.Args[1]
	tt, err := truetype.Parse(fonts.OpenSans_Regular_tff)
	if err != nil {
		log.Fatal(err)
	}

	lattice = entangler.NewLattice(context.TODO(), 3, 5, 5, 1015)
	//lattice = entangler.NewLattice(context.TODO(), 3, 5, 5, 259)
	lattice.RunInit()
	w, h := ebiten.ScreenSizeInFullscreen()
	windowXSize, windowYSize = w-10, int(float64(h)*0.5)
	displayXSize = (float64(xOffset) + (float64(lattice.NumDataBlocks/lattice.HorizontalStrands)+float64(0.5))*xSpace)
	zoom = float64(windowXSize) / displayXSize
	fmt.Printf("X: %v, Screen: %v, Zoom: %v\n", displayXSize, windowXSize, zoom)

	if _, err := os.Stat(logfile); os.IsNotExist(err) {
		fmt.Printf("Did not find logfile %v! Err: %v\n", logfile, err)
		os.Exit(1)
	}

	logParser = NewLogParser(logfile)
	if err := logParser.ParseLog(); err != nil {
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

func keyPressed(img *ebiten.Image, key ebiten.Key, presses int) {
	keyLock.Lock()
	defer keyLock.Unlock()
	if presses < keyPresses {
		return
	}

	if key.String() == ebiten.KeyLeft.String() {
		logParser.BlockCursor--
		for i := 0; i < len(lattice.Blocks); i++ {
			lattice.Blocks[i].Data = nil
			lattice.Blocks[i].IsUnavailable = false
			lattice.Blocks[i].WasDownloaded = false
		}
	} else if key.String() == ebiten.KeyRight.String() {
		logParser.BlockCursor++
	} else if key.String() == ebiten.KeyUp.String() {
		logParser.BlockCursor = logParser.BlockCursor + 10
	} else if key.String() == ebiten.KeyDown.String() {
		logParser.BlockCursor = logParser.BlockCursor - 10
		for i := 0; i < len(lattice.Blocks); i++ {
			lattice.Blocks[i].Data = nil
			lattice.Blocks[i].IsUnavailable = false
			lattice.Blocks[i].WasDownloaded = false
		}
	} else if key.String() == ebiten.KeyQ.String() {
		logParser.BlockCursor = len(logParser.TotalEntry[logParser.TotalCursor].BlockEntries)
	} else if key.String() == ebiten.Key1.String() {
		logParser.BlockCursor = 0
		logParser.TotalCursor--
		for i := 0; i < len(lattice.Blocks); i++ {
			lattice.Blocks[i].Data = nil
			lattice.Blocks[i].IsUnavailable = false
			lattice.Blocks[i].WasDownloaded = false
		}
	} else if key.String() == ebiten.Key2.String() {
		logParser.BlockCursor = 0
		logParser.TotalCursor++
		for i := 0; i < len(lattice.Blocks); i++ {
			lattice.Blocks[i].Data = nil
			lattice.Blocks[i].IsUnavailable = false
			lattice.Blocks[i].WasDownloaded = false
		}
	} else if key.String() == ebiten.Key3.String() {
		zoom -= zoominc
	} else if key.String() == ebiten.Key4.String() {
		zoom += zoominc
	} else if key.String() == ebiten.Key5.String() {
		columnOffset += columninc
	} else if key.String() == ebiten.Key6.String() {
		newColumn := columnOffset - columninc
		columnOffset = max(0, newColumn)
	}

	logParser.TotalCursor = min(max(0, logParser.TotalCursor), len(logParser.TotalEntry)-1)
	logParser.BlockCursor = min(max(0, logParser.BlockCursor), len(logParser.TotalEntry[logParser.TotalCursor].BlockEntries))

	logParser.ReadLog(lattice)
	if key.String() == ebiten.KeyH.String() {
		printHelp()
	}
	time.Sleep(300 * time.Millisecond)
	keyPresses++
}

func update(screen *ebiten.Image) error {
	zooming := false
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			go keyPressed(screen, k, keyPresses)
			if k.String() == ebiten.Key3.String() || k.String() == ebiten.Key4.String() {
				zooming = true
			}

			break
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	if zooming {
		ebiten.SetScreenScale(zoom)
		ebiten.SetScreenSize(int(float64(windowXSize)/zoom), windowYSize)
	}

	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	addCircle(screen, 400, 450, 25, color.Black, color.RGBA{0xff, 0, 0, 0xff})
	text.Draw(screen, "Unavailable", dataFont, 430, 462, color.Black)

	addCircle(screen, 650, 450, 25, color.Black, color.RGBA{0x33, 0x99, 0xff, 0xff})
	text.Draw(screen, "Repaired", dataFont, 680, 462, color.Black)

	addCircle(screen, 900, 450, 25, color.Black, color.RGBA{0x0, 0xff, 0, 0xff})
	text.Draw(screen, "Downloaded", dataFont, 930, 462, color.Black)

	text.Draw(screen, fmt.Sprintf("Entry: %d / %d", logParser.TotalCursor+1,
		len(logParser.TotalEntry)),
		dataFont, 1380, 462, color.Black)

	text.Draw(screen, fmt.Sprintf("Event: %d / %d", logParser.BlockCursor,
		len(logParser.TotalEntry[logParser.TotalCursor].BlockEntries)),
		dataFont, 1130, 462, color.Black)

	numBlocks := len(lattice.Blocks)

	for i := (columnOffset * 5); i < (numBlocks + (columnOffset * 5)); i++ {
		block := lattice.Blocks[i%numBlocks]
		if !block.IsParity || (!block.HasData() && !block.IsUnavailable) {
			continue
		}
		var leftPos, rightPos int

		if len(block.Left) == 0 || block.Left[0].Position < 1 {
			rightPos = block.Right[0].Position + lattice.NumDataBlocks
			r, h, l := entangler.GetBackwardNeighbours(rightPos, lattice.S, lattice.P)
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

		var clr color.Color
		if !block.HasData() && block.IsUnavailable {
			clr = color.RGBA{0xff, 0, 0, 0xff} // Dotted red line
		} else if !block.WasDownloaded {
			clr = color.RGBA{0x33, 0x99, 0xff, 0xff} // Blue line
		} else {
			clr = color.RGBA{0x0, 0xff, 0, 0xff} // Green line
		}
		// switch block.Class {
		// case entangler.Horizontal:
		// 	clr = color.RGBA{0, 0xff, 0, 0xff}
		// case entangler.Right:
		// 	clr = color.RGBA{0, 0, 0xff, 0xff}
		// case entangler.Left:
		// 	clr = color.Black
		// }
		co := -columnOffset
		if i >= numBlocks {
			co += lattice.NumDataBlocks/lattice.HorizontalStrands + 1
		}

		addParityBetweenDatablock(screen, leftPos, rightPos, clr, 8, co, lattice.HorizontalStrands, block.Class)
	}
	for i := (columnOffset * 5); i < (numBlocks + (columnOffset * 5)); i++ {
		bl := lattice.Blocks[i%numBlocks]
		if bl.IsParity {
			continue
		}
		var clr color.Color
		if bl.HasData() {
			if !bl.WasDownloaded {
				clr = color.RGBA{0x33, 0x99, 0xff, 0xff}
			} else {
				clr = color.RGBA{0x0, 0xff, 0, 0xff}
			}
		} else if bl.IsUnavailable {
			clr = color.RGBA{0xff, 0x0, 0x0, 0xff}
		} else {
			clr = color.RGBA{0xc8, 0xc8, 0xc8, 0xff}
		}
		co := columnOffset
		if i >= numBlocks {
			co -= lattice.NumDataBlocks/lattice.HorizontalStrands + 1
		}
		addDataBlock(screen, dataRadius, color.Black,
			clr, color.Black,
			lattice.Blocks[i%numBlocks].Position, co, lattice.HorizontalStrands)
	}
	return nil
}

func main() {
	ebiten.SetMaxTPS(60)
	ebiten.SetRunnableInBackground(true)
	printHelp()
	if err := ebiten.Run(update, int(float64(windowXSize)/zoom), windowYSize, zoom, windowTitle); err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	fmt.Printf("--- Help for Entangle Visualizer ---\n" +
		" [h]       - This help.\n" +
		"--- Navigation ---\n" +
		" [q]       - Scroll to the end of the log.\n" +
		" [R ARROW] - Next event in the log.\n" +
		" [L ARROW] - Previous event in the log.\n" +
		" [U ARROW] - Next 10 events in the log.\n" +
		" [D ARROW] - Previous 10 events in the log.\n" +
		" [1]       - Jump to the previous log entry.\n" +
		" [2]       - Jump to the next log entry.\n" +
		"--- Window ---\n" +
		" [3]       - Zoom out.\n" +
		" [4]       - Zoom in.\n" +
		" [5]       - Scroll entire window to the left.\n" +
		" [6]       - Scroll entire window to the right.\n",
	)
}
