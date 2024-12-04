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
var ListenAddr *string
var Headless *bool
var WatchDir *string  // WatchDir is the directory to watch for new files
var OutputDir *string // OutputDir is the directory to output converted files
var UploadDir *string // UploadDir is the directory to store files being uploaded by the user

/**
 * CONVERSION OPTIONS
 **/
var conv chan int        // conv is a channel to limit the number of concurrent conversions
var blockAuto chan bool  // blockAuto is a channel to block automatic conversion
var maxConcurrent *int   // max number of concurrent conversions
var AutoConvert *bool    // true to automatically convert files in the watch directory
var DeleteAfter *bool    // true to delete source files after conversion
var IgnoreExisting *bool // true to overwrite already converted files

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
	opts.HelpCommand("help", opts.Alias("h"))
	Port = opts.Int("port", 8080, opts.Description("Port to listen on"), opts.Alias("p"))
	ListenAddr = opts.String("listen", "127.0.0.1", opts.Description("Address to listen on"), opts.Alias("l"))
	Headless = opts.Bool("headless", false, opts.Description("Run in headless mode"), opts.Alias("H"))

	maxConcurrent = opts.Int("concurrency", 1, opts.Description("Max number of concurrent conversions"), opts.Alias("c"))
	maxQueueSize = opts.Int("queue-size", 100, opts.Description("Max number of files to queue"), opts.Alias("q"))
	WatchDir = opts.String("watch-dir", "./media/input", opts.Description("Directory to watch for new files"), opts.Alias("w"))
	OutputDir = opts.String("output-dir", "./media/output", opts.Description("Directory to output converted files"), opts.Alias("o"))
	UploadDir = opts.String("upload-dir", "./media/upload", opts.Description("Directory to store files being uploaded by the user"), opts.Alias("u"))

	AutoConvert = opts.Bool("auto-convert", true, opts.Description("Automatically convert files in the watch directory"), opts.Alias("a"))
	DeleteAfter = opts.Bool("delete-after", false, opts.Description("Delete source files after conversion"), opts.Alias("d"))
	IgnoreExisting = opts.Bool("ignore-existing", false, opts.Description("Overwrite already converted files"), opts.Alias("i"))

	opts.Parse(os.Args[1:])
	if opts.Called("help") {
		fmt.Fprintf(os.Stdout, opts.Help())
		os.Exit(0)
	}

	if opts.Called("listen") == true && *ListenAddr == "127.0.0.1" {
		*ListenAddr = "0.0.0.0"
	}

	fmt.Println("Will listen on", *ListenAddr)
	fmt.Println("watching", *WatchDir)
	fmt.Println("outputting to", *OutputDir)
	fmt.Println("uploading to", *UploadDir)
	conv = make(chan int, *maxConcurrent)
	fileQueue = make(chan File, *maxQueueSize)
}
