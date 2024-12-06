package actions

import "github.com/asishshaji/notion-backup/models"

// data that is shared across the actions
// Think of this as the request that goes through all the chain of actions
type SharedData struct {
	ExportType         models.ExportType
	TaskId             string
	ExportURL          string
	DownloadedFilePath string
}

type Action interface {
	// the handler for the action
	Act(*SharedData) error
	// the name of the action
	String() string
}
