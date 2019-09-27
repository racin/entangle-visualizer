package main

// fmt.Printf("%t,%d,%d,%d,%t,%d,%d\n", block.IsParity, block.Position, block.LeftPos(0), block.RightPos(0), block.HasData(), start, time.Now().UnixNano())
type ParityEntry struct {
	IsParity  bool
	Position  int
	LeftPos   int
	RightPos  int
	HasData   bool
	StartTime int64
	EndTime   int64
	Error     string // HTTP 404
}

type TotalEntry struct {
	Datablocks   int
	Parityblocks int
	StartTime    int64
	EndTime      int64
	Error        string // Timeout
}
