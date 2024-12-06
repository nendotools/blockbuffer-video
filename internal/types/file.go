package types

type FileStatus string

const (
	New             FileStatus = "new"
	Queued          FileStatus = "queued"
	Processing      FileStatus = "processing"
	Completed       FileStatus = "completed"
	CompleteDeleted FileStatus = "completed-deleted"
	Cancelled       FileStatus = "cancelled"
	Rejected        FileStatus = "rejected"
	Failed          FileStatus = "failed"
	Deleted         FileStatus = "deleted"
)

type File struct {
	ID       string     `json:"id"`
	FilePath string     `json:"filePath"`
	Status   FileStatus `json:"status"`
	Progress float32    `json:"progress"`
	Duration float64    `json:"duration"`
}
