package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/racin/snarl/entangler"
	"github.com/racin/snarl/swarmconnector"
)

// fmt.Printf("%t,%d,%d,%d,%t,%d,%d\n", block.IsParity, block.Position, block.LeftPos(0), block.RightPos(0), block.HasData(), start, time.Now().UnixNano())

type timePeriod struct {
	StartTime int64
	EndTime   int64
}

type BlockEntry struct {
	IsParity       bool
	Position       int
	LeftPos        int
	RightPos       int
	HasData        bool
	DownloadTime   timePeriod
	RepairTime     timePeriod
	Error          string // HTTP 404
	DownloadStatus DownloadStatus
	RepairStatus   RepairStatus
}

func NewBlockEntryString(entry []string, err string) BlockEntry {
	IsParity, _ := strconv.ParseBool(entry[0])
	Position, _ := strconv.Atoi(entry[1])
	LeftPos, _ := strconv.Atoi(entry[2])
	RightPos, _ := strconv.Atoi(entry[3])
	HasData, _ := strconv.ParseBool(entry[4])
	DownloadStartTime, _ := strconv.ParseInt(entry[5], 10, 64)
	DownloadEndTime, _ := strconv.ParseInt(entry[6], 10, 64)
	var RepairStartTime, RepairEndTime int64
	var dlStatus DownloadStatus
	var repStatus RepairStatus
	if len(entry) == 8 {
		WasDownloaded, _ := strconv.ParseBool(entry[7][:len(entry[7])-1])
		if WasDownloaded {
			dlStatus = DownloadSuccess
		} else {
			dlStatus = DownloadFailed
		}
	} else if len(entry) == 11 {
		dlStatus = ConvertDLStatus(entry[7])
		RepairStartTime, _ = strconv.ParseInt(entry[8], 10, 64)
		RepairEndTime, _ = strconv.ParseInt(entry[9], 10, 64)
		repStatus = ConvertRepStatus(entry[10][:len(entry[10])-1])
	}

	return NewBlockEntry(IsParity, Position, LeftPos, RightPos,
		HasData, DownloadStartTime, DownloadEndTime, RepairStartTime,
		RepairEndTime, err, dlStatus, repStatus)
}

func NewBlockEntry(IsParity bool, Position, LeftPos, RightPos int,
	HasData bool, DownloadStartTime, DownloadEndTime, RepairStartTime, RepairEndTime int64,
	err string, downloadStatus DownloadStatus,
	repairStatus RepairStatus) BlockEntry {
	return BlockEntry{IsParity: IsParity, Position: Position,
		LeftPos: LeftPos, RightPos: RightPos, HasData: HasData,
		DownloadTime: timePeriod{StartTime: DownloadStartTime, EndTime: DownloadEndTime},
		RepairTime:   timePeriod{StartTime: RepairStartTime, EndTime: RepairEndTime},
		Error:        err, DownloadStatus: downloadStatus, RepairStatus: repairStatus}
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

	var getSortTime func(BlockEntry) int64 = func(entry BlockEntry) int64 {
		if entry.RepairStatus == RepairSuccess || entry.RepairStatus == RepairFailed {
			return entry.RepairTime.EndTime
		} else if entry.RepairStatus == RepairPending {
			return entry.RepairTime.StartTime
		} else if entry.DownloadStatus == DownloadSuccess || entry.DownloadStatus == DownloadFailed {
			return entry.DownloadTime.EndTime
		} else {
			return entry.DownloadTime.StartTime
		}
	}
	sort.Slice(BlockEntries, func(i, j int) bool {
		return getSortTime(BlockEntries[i]) < getSortTime(BlockEntries[j])
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

func (l *LogParser) ParseLog(offset ...int) error {
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
		entry := strings.Split(line, ",")
		switch len(entry) {
		case 11:
			fallthrough
		case 8:
			newEntry := NewBlockEntryString(entry, blockError)
			if newEntry.RepairStatus == RepairPending || newEntry.DownloadStatus == DownloadPending {
				continue
			}
			if len(offset) == 2 && offset[0] != 0 && offset[1] != 0 {
				if newEntry.Position <= offset[0] || newEntry.Position > offset[1] {
					continue
				}
				newEntry.Position -= offset[0]
				newEntry.LeftPos -= offset[0]
				newEntry.RightPos -= offset[0]
			}
			blockEntries = append(blockEntries, newEntry)
			blockError = ""
		case 2:
			if strings.Contains(line, "Downloaded total") {
				db := strings.Split(entry[0], ": ")
				pb := strings.Split(entry[1], ": ")
				Datablocks, _ = strconv.Atoi(db[1])
				Parityblocks, _ = strconv.Atoi(pb[1][:len(pb[1])-1])
			} else {
				StartTime, _ := strconv.ParseInt(entry[0], 10, 64)
				EndTime, _ := strconv.ParseInt(entry[1][:len(entry[1])-1], 10, 64)
				if len(blockEntries) == 0 {
					continue
				}
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
	for i := 0; i < l.BlockCursor; i++ {
		fmt.Printf("Reading log. Cursor: %d", i)
		entry := l.TotalEntry[l.TotalCursor].BlockEntries[i]
		if entry.IsParity {
			fmt.Printf(" - Drawing Parity. Left: %d, Right: %d\n", entry.LeftPos, entry.RightPos)
			if entry.RightPos > lattice.NumDataBlocks {
				continue
			} else if entry.LeftPos < 1 {
				rightData := lattice.Blocks[entry.RightPos-1]
				var j int
				r, h, l := entangler.GetBackwardNeighbours(entry.RightPos, lattice.S, lattice.P)
				if entry.LeftPos == r {
					j = 1
				} else if entry.LeftPos == h {
					j = 0
				} else if entry.LeftPos == l {
					j = 2
				}
				rightData.Left[j].WasDownloaded = entry.DownloadStatus == DownloadSuccess
				if !entry.HasData || entry.Error != "" {
					rightData.Left[j].IsUnavailable = true
					fmt.Printf(" - Parity Unavailable. Left: %d, Right: %d\n", entry.LeftPos, entry.RightPos)
				} else {
					rightData.Left[j].Data = make([]byte, swarmconnector.ChunkSizeOffset+1)
				}
			} else {
				// TODO: Increment with numdatablocks .
				leftData := lattice.Blocks[entry.LeftPos-1]
				for j := 0; j < len(leftData.Right); j++ {
					rightData := leftData.Right[j].Right[0]
					if len(leftData.Right[j].Right) > 0 &&
						rightData.Position == entry.RightPos {

						leftData.Right[j].WasDownloaded = entry.DownloadStatus == DownloadSuccess
						if !entry.HasData || entry.Error != "" {
							leftData.Right[j].IsUnavailable = true
							fmt.Printf(" - Parity Unavailable. Left: %d, Right: %d\n", entry.LeftPos, entry.RightPos)
						} else {
							leftData.Right[j].Data = make([]byte, swarmconnector.ChunkSizeOffset+1)
						}
					}
				}
			}
		} else {
			b := lattice.Blocks[entry.Position-1]
			if entry.DownloadStatus == DownloadFailed && !entry.HasData {
				fmt.Printf(" - Data Unavailable. Position:%d\n", entry.Position)
				b.IsUnavailable = true
			} else if (entry.DownloadStatus == DownloadFailed || entry.RepairStatus == RepairSuccess) && entry.HasData {
				fmt.Printf(" - Data Reconstructed. Position:%d\n", entry.Position)
				fmt.Printf("Entry: %v\n", entry)
				b.Data = make([]byte, swarmconnector.ChunkSizeOffset+1)
				b.WasDownloaded = entry.DownloadStatus == DownloadSuccess
			} else if entry.DownloadStatus == DownloadPending {
				fmt.Printf(" - Data pending. Position:%d\n", entry.Position)
			} else if entry.DownloadStatus == DownloadSuccess {
				fmt.Printf(" - Data success. Position:%d\n", entry.Position)
				b.Data = make([]byte, swarmconnector.ChunkSizeOffset+1)
				b.WasDownloaded = entry.DownloadStatus == DownloadSuccess
			} else {
				fmt.Printf("%v\n", entry)
			}
		}
	}
}
