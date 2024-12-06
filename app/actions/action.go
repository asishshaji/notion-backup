package actions

import "github.com/asishshaji/notion-backup/models"

type SharedData struct {
	ExportType models.ExportType
	TaskId     string
	ExportURL  string
}

type Action interface {
	Act(*SharedData) error
}
