package internal

import (
	"fmt"
	"os"
	"time"

	getopts "github.com/DavidGamba/go-getoptions"
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
var Port *int
var WatchDir *string

var conv = make(chan int, maxConcurrent)

func init() {
	opts := getopts.New()
	opts.HelpCommand("h", opts.Alias("help"))
	Port = opts.Int("p", 8080, opts.Description("Port to listen on"), opts.Alias("port"))
	WatchDir = opts.String("w", "/mnt/internal-ssd/auto-convert/input", opts.Description("Directory to watch for new files"), opts.Alias("watch-dir"))

	opts.Parse(os.Args[1:])
	if opts.Called("help") {
		fmt.Fprintf(os.Stdout, opts.Help())
		os.Exit(0)
	}

	// print the Port and WatchDir
	fmt.Println("Port: ", *Port)
	fmt.Println("WatchDir: ", *WatchDir)
}
