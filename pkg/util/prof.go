package util

import (
	"fmt"
	"log/slog"
	"runtime"
)

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

type ByteSize float64

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

func PrintResourceUsage() {
	printNumCpu()
	printMemUsage()
}

func printNumCpu() {
	slog.Info("cpu usage", "NumCPU", runtime.NumCPU())
}

func printMemUsage() {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	slog.Info("memory usage",
		"Alloc", ByteSize(m.Alloc),
		"TotalAlloc", ByteSize(m.TotalAlloc),
		"Sys", ByteSize(m.Sys),
		"NumGC", m.NumGC,
	)
	//log.Printf("Alloc = %v", ByteSize(m.Alloc))
	//log.Printf("TotalAlloc = %v", ByteSize(m.TotalAlloc))
	//log.Printf("Sys = %v", ByteSize(m.Sys))
	//log.Printf("NumGC = %v\n", m.NumGC)
}
