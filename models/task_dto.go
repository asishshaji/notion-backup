package models

type ExportOptions struct {
	ExportType string `json:"exportType"`
	TimeZone   string `json:"timeZone"`
	Locale     string `json:"locale"`
}
type TaskRequest struct {
	SpaceId              string        `json:"spaceId"`
	ExportOptions        ExportOptions `json:"exportOptions"`
	ShouldExportComments bool          `json:"shouldExportComments"`
}
type Task struct {
	EventName string      `json:"eventName"`
	Request   TaskRequest `json:"request"`
}

type CreateTaskDTO struct {
	T Task `json:"task"`
}

type CreateTaskResponseDTO struct {
	TaskId string `json:"taskId"`
}
