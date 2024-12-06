package constants

import "github.com/asishshaji/notion-backup/models"

const (
	MardownExportType models.ExportType = "markdown"
	HtmlExportType    models.ExportType = "html"
)

const NOTION_API_BASE_URL = "https://www.notion.so/api/v3"
const NOTION_API_ENQUEUE_URL = "https://www.notion.so/api/v3/enqueueTask"
const NOTION_API_GET_TASKS_URL = "https://www.notion.so/api/v3/getTasks"
