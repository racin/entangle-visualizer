package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

// fmt.Printf("%t,%d,%d,%d,%t,%d,%d\n", block.IsParity, block.Position, block.LeftPos(0), block.RightPos(0), block.HasData(), start, time.Now().UnixNano())

type BlockEntry struct {
	IsParity  bool
	Position  int
	LeftPos   int
	RightPos  int
	HasData   bool
	StartTime int64
	EndTime   int64
	Error     string // HTTP 404
}

func NewBlockEntryString(entry []string, err string) BlockEntry {
	IsParity, _ := strconv.ParseBool(entry[0])
	Position, _ := strconv.Atoi(entry[1])
	LeftPos, _ := strconv.Atoi(entry[2])
	RightPos, _ := strconv.Atoi(entry[3])
	HasData, _ := strconv.ParseBool(entry[4])
	StartTime, _ := strconv.ParseInt(entry[5], 10, 64)
	EndTime, _ := strconv.ParseInt(entry[6][:len(entry[6])-2], 10, 64)
	return NewBlockEntry(IsParity, Position, LeftPos, RightPos,
		HasData, StartTime, EndTime, err)
}

func NewBlockEntry(IsParity bool, Position, LeftPos, RightPos int,
	HasData bool, StartTime, EndTime int64, err string) BlockEntry {
	return BlockEntry{IsParity: IsParity, Position: Position,
		LeftPos: LeftPos, RightPos: RightPos, HasData: HasData,
		StartTime: StartTime, EndTime: EndTime, Error: err}
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
	return TotalEntry{Datablocks: Datablocks, Parityblocks: Parityblocks,
		StartTime: StartTime, EndTime: EndTime, Error: err,
		BlockEntries: BlockEntries}
}

type LogParser struct {
	Path       string
	Cursor     int
	TotalEntry []TotalEntry
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

func (l *LogParser) ReadLog() error {
	// Process each line.
	file, err := os.Open(l.Path)
	defer file.Close()

	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	var blockEntrySize float64 = 0
	blockEntries := make([]BlockEntry, int(blockEntrySize))
	var blockError, totalError = "", ""

	var line string
	var Datablocks, Parityblocks = 0, 0
	for {
		line, err = reader.ReadString('\n')

		switch strings.Count(line, ",") {
		case 6:
			fmt.Println("BlockEntry")
			entry := strings.Split(line, ",")
			blockEntries = append(blockEntries, NewBlockEntryString(entry, blockError))
			blockError = ""
		case 1:
			entry := strings.Split(line, ",")
			if strings.Contains(line, "Downloaded total") {
				db := strings.Split(entry[0], ": ")
				pb := strings.Split(entry[1], ": ")
				Datablocks, _ = strconv.Atoi(db[1])
				Parityblocks, _ = strconv.Atoi(pb[1][:len(pb[1])-2])
			} else {
				StartTime, _ := strconv.ParseInt(entry[0], 10, 64)
				EndTime, _ := strconv.ParseInt(entry[1][:len(entry[1])-2], 10, 64)
				l.TotalEntry = append(l.TotalEntry, NewTotalEntry(Datablocks, Parityblocks, StartTime,
					EndTime, totalError, blockEntries))
				blockEntrySize = math.Max(float64(blockEntrySize), float64(Datablocks+Parityblocks))
				blockEntries = make([]BlockEntry, int(blockEntrySize))
				Datablocks, Parityblocks = 0, 0
				totalError, blockError = "", ""
			}
		default:
			if strings.Contains(line, "404 Not Found") {
				blockError = line
			} else {
				totalError += line
			}
		}
		// Process the line here.
		fmt.Println(" > > " + limitLength(line, 50))

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		fmt.Printf(" > Failed!: %v\n", err)
	}

	return nil

	// If currently not on a TotalStruct, create one. "Downloaded file" finishes

	// After TotalStruct is finished. Sort each BlockEntry

}
