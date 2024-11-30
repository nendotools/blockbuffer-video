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

type conversionDone string
type progressMsg struct {
	percent    float64
	conversion string
}

var Port *int
var WatchDir *string  // WatchDir is the directory to watch for new files
var OutputDir *string // OutputDir is the directory to output converted files
var UploadDir *string // UploadDir is the directory to store files being uploaded by the user

/**
 * CONVERSION OPTIONS
 **/
var conv chan int      // conv is a channel to limit the number of concurrent conversions
var maxConcurrent *int // max number of concurrent conversions

/**
*  FILE QUEUE OPTIONS
 **/
var maxQueueSize *int
var fileQueue chan File // fileQueue is a channel to queue files to be processed
var skipList = make(map[string]bool)

const maxCheckInterval = 5 * time.Second
const maxCheckRepeat = 30 // 5 minutes, to support larger files or slow writes
const maxQueueRetry = 3   // failed files will be retried up to 3 times

func init() {
	opts := getopts.New()
	opts.HelpCommand("h", opts.Alias("help"))
	Port = opts.Int("p", 8080, opts.Description("Port to listen on"), opts.Alias("port"))
	WatchDir = opts.String("w", "/mnt/internal-ssd/auto-convert/input", opts.Description("Directory to watch for new files"), opts.Alias("watch-dir"))
	OutputDir = opts.String("o", "/mnt/internal-ssd/auto-convert/output", opts.Description("Directory to output converted files"), opts.Alias("output-dir"))
	UploadDir = opts.String("u", "/mnt/internal-ssd/auto-convert/upload", opts.Description("Directory to store files being uploaded by the user"), opts.Alias("upload-dir"))
	maxConcurrent = opts.Int("c", 8, opts.Description("Max number of concurrent conversions"), opts.Alias("concurrency"))
	maxQueueSize = opts.Int("q", 100, opts.Description("Max number of files to queue"), opts.Alias("queue-size"))

	opts.Parse(os.Args[1:])
	if opts.Called("help") {
		fmt.Fprintf(os.Stdout, opts.Help())
		os.Exit(0)
	}

	conv = make(chan int, *maxConcurrent)
	fileQueue = make(chan File, *maxQueueSize)

	// print the Port and WatchDir
	fmt.Println("Port: ", *Port)
	fmt.Println("WatchDir: ", *WatchDir)
}
