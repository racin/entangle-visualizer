package main

import (
	"bufio"
	"fmt"
	"github.com/racin/entangle/entangler"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// fmt.Printf("%t,%d,%d,%d,%t,%d,%d\n", block.IsParity, block.Position, block.LeftPos(0), block.RightPos(0), block.HasData(), start, time.Now().UnixNano())

type BlockEntry struct {
	IsParity      bool
	Position      int
	LeftPos       int
	RightPos      int
	HasData       bool
	StartTime     int64
	EndTime       int64
	Error         string // HTTP 404
	WasDownloaded bool
}

func NewBlockEntryString(entry []string, err string) BlockEntry {
	IsParity, _ := strconv.ParseBool(entry[0])
	Position, _ := strconv.Atoi(entry[1])
	LeftPos, _ := strconv.Atoi(entry[2])
	RightPos, _ := strconv.Atoi(entry[3])
	HasData, _ := strconv.ParseBool(entry[4])
	StartTime, _ := strconv.ParseInt(entry[5], 10, 64)
	EndTime, _ := strconv.ParseInt(entry[6], 10, 64)
	WasDownloaded, _ := strconv.ParseBool(entry[7][:len(entry[7])-1])
	return NewBlockEntry(IsParity, WasDownloaded, Position, LeftPos, RightPos,
		HasData, StartTime, EndTime, err)
}

func NewBlockEntry(IsParity, WasDownloaded bool, Position, LeftPos, RightPos int,
	HasData bool, StartTime, EndTime int64, err string) BlockEntry {
	return BlockEntry{IsParity: IsParity, Position: Position,
		LeftPos: LeftPos, RightPos: RightPos, HasData: HasData,
		StartTime: StartTime, EndTime: EndTime, Error: err,
		WasDownloaded: WasDownloaded}
}

type TotalEntry struct {
	Datablocks   int
	Parityblocks int
	StartTime    int64
	EndTime      int64
	Error        string // Timeout
	BlockEntries []BlockEntry
}

func NewTotalEntry(Datablocks, Parityblocks int, StartTime,
	EndTime int64, err string, BlockEntries []BlockEntry) TotalEntry {

	sort.Slice(BlockEntries, func(i, j int) bool {
		return BlockEntries[i].StartTime < BlockEntries[j].StartTime
	})

	return TotalEntry{Datablocks: Datablocks, Parityblocks: Parityblocks,
		StartTime: StartTime, EndTime: EndTime, Error: err,
		BlockEntries: BlockEntries}
}

type LogParser struct {
	Path        string
	TotalCursor int
	BlockCursor int
	TotalEntry  []TotalEntry
}

func NewLogParser(filepath string) *LogParser {
	return &LogParser{
		Path:       filepath,
		TotalEntry: make([]TotalEntry, 0, 1),
	}
}

func limitLength(s string, length int) string {
	if len(s) < length {
		return s
	}
	return s[:length]
}

func (l *LogParser) ParseLog() error {
	file, err := os.Open(l.Path)
	defer file.Close()

	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	var blockEntrySize int = 0
	blockEntries := make([]BlockEntry, int(blockEntrySize))
	var blockError, totalError = "", ""

	var line string
	var Datablocks, Parityblocks = 0, 0
	for {
		line, err = reader.ReadString('\n')

		switch strings.Count(line, ",") {
		case 7:
			entry := strings.Split(line, ",")
			blockEntries = append(blockEntries, NewBlockEntryString(entry, blockError))
			blockError = ""
		case 1:
			entry := strings.Split(line, ",")
			if strings.Contains(line, "Downloaded total") {
				db := strings.Split(entry[0], ": ")
				pb := strings.Split(entry[1], ": ")
				Datablocks, _ = strconv.Atoi(db[1])
				Parityblocks, _ = strconv.Atoi(pb[1][:len(pb[1])-1])
			} else {
				StartTime, _ := strconv.ParseInt(entry[0], 10, 64)
				EndTime, _ := strconv.ParseInt(entry[1][:len(entry[1])-1], 10, 64)
				l.TotalEntry = append(l.TotalEntry, NewTotalEntry(Datablocks, Parityblocks, StartTime,
					EndTime, totalError, blockEntries))
				blockEntrySize = max(blockEntrySize, Datablocks+Parityblocks)
				blockEntries = make([]BlockEntry, 0, blockEntrySize)
				Datablocks, Parityblocks = 0, 0
				totalError, blockError = "", ""
			}
		default:
			if strings.Contains(line, "FATAL ERROR. STOPPING DOWNLOAD.") {
				blockEntries = make([]BlockEntry, 0, blockEntrySize)
				Datablocks, Parityblocks = 0, 0
				totalError, blockError = "", ""
			} else if strings.Contains(line, "404 Not Found") {
				blockError = line
			} else {
				totalError += line
			}
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		fmt.Printf(" > Failed!: %v\n", err)
	}

	return nil
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
func (l *LogParser) ReadLog(lattice *entangler.Lattice) {
	for i := 0; i < min(max(0, l.BlockCursor), len(l.TotalEntry[l.TotalCursor].BlockEntries)); i++ {
		fmt.Printf("Reading log. Cursor: %d", i)
		entry := l.TotalEntry[l.TotalCursor].BlockEntries[i]
		if entry.IsParity {
			fmt.Printf(" - Drawing Parity. Left: %d, Right: %d\n", entry.LeftPos, entry.RightPos)
			if entry.LeftPos < 1 || entry.RightPos > lattice.NumDataBlocks {
				continue
			}
			leftData := lattice.Blocks[entry.LeftPos-1]
			for j := 0; j < len(leftData.Right); j++ {
				rightData := leftData.Right[j].Right[0]
				if len(leftData.Right[j].Right) > 0 &&
					rightData.Position == entry.RightPos {
					leftData.Right[j].Data = make([]byte, 1)
					leftData.Right[j].WasDownloaded = entry.WasDownloaded
					if !entry.HasData || entry.Error != "" {
						leftData.Right[j].IsUnavailable = true
					} else {
						// Check if data was reconstructed.
						if len(leftData.Left) > 0 && leftData.Left[j].HasData() && !leftData.HasData() {
							leftData.Data = make([]byte, 1)
							leftData.WasDownloaded = false
						}
						if len(rightData.Right) > 0 && rightData.Right[j].HasData() && !rightData.HasData() {
							rightData.Data = make([]byte, 1)
							rightData.WasDownloaded = false
						}

						// Check if Parity was reconstructed
						if len(leftData.Left) > 0 && !leftData.Left[j].HasData() && leftData.HasData() {
							leftData.Left[j].Data = make([]byte, 1)
							leftData.Left[j].WasDownloaded = false
							leftData.Left[j].IsUnavailable = false
						}
						if len(rightData.Right) > 0 && !rightData.Right[j].HasData() && rightData.HasData() {
							rightData.Right[j].Data = make([]byte, 1)
							rightData.Right[j].WasDownloaded = false
							rightData.Right[j].IsUnavailable = false
						}
					}
				}
			}
		} else {
			b := lattice.Blocks[entry.Position-1]
			if !entry.HasData || entry.Error != "" {
				fmt.Printf(" - Data Unavailable. Position:%d\n", entry.Position)
				b.IsUnavailable = true
			} else {
				fmt.Printf(" - Data Reconstructed. Position:%d\n", entry.Position)
				b.Data = make([]byte, 1)
				b.WasDownloaded = entry.WasDownloaded
			}
		}
	}
}
