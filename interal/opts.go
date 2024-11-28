package internal

import (
	"flag"
	"time"
)

// fileQueue is a channel to queue files to be processed
// This is used to avoid processing the same file multiple times
// and prevent trigger multiple conversions at once
var fileQueue = make(chan File, 100)
var skipList = make(map[string]bool)

const OutputDir = "/mnt/internal-ssd/auto-convert/output"

type conversionDone string
type progressMsg struct {
	percent    float64
	conversion string
}

const maxQueueSize = 100
const maxCheckInterval = 5 * time.Second
const maxCheckRepeat = 30 // 5 minutes, to support larger files or slow writes
const maxQueueRetry = 3   // failed files will be retried up to 3 times
const maxConcurrent = 8   // max number of concurrent conversions

var conv = make(chan int, maxConcurrent)

var (
	Port     = flag.String("port", "8080", "Port to listen on")
	WatchDir = flag.String("watchDir", "/mnt/internal-ssd/auto-convert/input", "Directory to watch for files")
)

func init() {
	flag.Parse()
}
