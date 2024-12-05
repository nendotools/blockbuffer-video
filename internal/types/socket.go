package types

type MessageType string

const (
	RefreshFiles MessageType = "refresh_files"
	UpdateFile   MessageType = "update_file"
	CreateFile   MessageType = "create_file"
	DeleteFile   MessageType = "delete_file"
)

type Message struct {
	MessageType MessageType `json:"type"`
	MustSend    bool        `json:"must_send"`
	Data        interface{} `json:"data"`
}
