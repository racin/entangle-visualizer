package main

import "fmt"

type DownloadStatus int

const (
	NoDownload      = iota // Did not attempt to download yet
	DownloadPending        // Download pending
	DownloadSuccess        // Download finished and HasData() = true
	DownloadFailed         // Download finished and HasData() = false
)

func (s DownloadStatus) String() string {
	switch s {
	case NoDownload:
		return "NoDownload"
	case DownloadPending:
		return "DownloadPending"
	case DownloadSuccess:
		return "DownloadSuccess"
	case DownloadFailed:
		return "DownloadFailed"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}

func ConvertDLStatus(dlStatus string) DownloadStatus {
	switch dlStatus {
	case "NoDownload":
		return NoDownload
	case "DownloadPending":
		return DownloadPending
	case "DownloadSuccess":
		return DownloadSuccess
	case "DownloadFailed":
		return DownloadFailed
	default:
		return NoDownload
	}
}

func ConvertRepStatus(repairStatus string) RepairStatus {
	switch repairStatus {
	case "NoRepair":
		return NoRepair
	case "RepairPending":
		return RepairPending
	case "RepairSuccess":
		return RepairSuccess
	case "RepairFailed":
		return RepairFailed
	default:
		return NoRepair
	}
}

type RepairStatus int

const (
	NoRepair      = iota // Did not attempt to repair
	RepairPending        // We started the repair process for this block.
	RepairSuccess        // HasData() = true [Download initially failed or was never attempted]
	RepairFailed         // HasData() = false [Download initially failed or was never attempted]
)

func (s RepairStatus) String() string {
	switch s {
	case NoRepair:
		return "NoRepair"
	case RepairPending:
		return "RepairPending"
	case RepairSuccess:
		return "RepairSuccess"
	case RepairFailed:
		return "RepairFailed"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}
