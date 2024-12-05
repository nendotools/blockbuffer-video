package store

import (
	"sync"

	opts "blockbuffer/internal/settings"
	types "blockbuffer/internal/types"
)

var FileListMutex = &sync.Mutex{}
var FileList = make(map[string]types.File) // FileList is a map of file ID to file
// var FileQueue chan types.File              // fileQueue is a channel to queue files to be processed
var FileQueue = make(chan types.File, *opts.MaxQueueSize)

func UpdateFile(file types.File) {
	FileListMutex.Lock()
	FileList[file.ID] = file
	FileListMutex.Unlock()
}
