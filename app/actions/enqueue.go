package actions

import (
	"encoding/json"

	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/constants"
	"github.com/asishshaji/notion-backup/models"
)

type EnqueueAction struct {
	HttpClient *httpclient.HTTPClient
}

func (enqueueAction EnqueueAction) Act(s *SharedData) error {
	taskRequest := models.CreateTaskDTO{
		T: models.Task{
			EventName: "exportSpace",
			Request: models.TaskRequest{
				SpaceId: "",
				ExportOptions: models.ExportOptions{
					ExportType: string(s.ExportType),
					TimeZone:   "America/New_York",
					Locale:     "en",
				},
				ShouldExportComments: false,
			},
		},
	}

	marshalledTaskRequest, err := json.Marshal(taskRequest)
	if err != nil {
		return err
	}

	resp, err := enqueueAction.HttpClient.Post(constants.NOTION_API_ENQUEUE_URL, marshalledTaskRequest)
	if err != nil {
		return err
	}
	var taskResp models.CreateTaskResponseDTO

	if err = json.Unmarshal(resp, &taskResp); err != nil {
		return err
	}

	// save the taskId to shared data
	s.TaskId = taskResp.TaskId

	return nil
}
