package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/constants"
	"github.com/asishshaji/notion-backup/models"
)

type EnqueueAction struct {
	HttpClient      *httpclient.HTTPClient
	NOTION_SPACE_ID string // the workspace to download
}

func (EnqueueAction) String() string {
	return "EnqueueAction"
}

func (enqueueAction EnqueueAction) createEnqueueRequest(exportType, spaceId string) (*http.Request, error) {
	taskRequest := models.CreateTaskDTO{
		T: models.Task{
			EventName: "exportSpace",
			Request: models.TaskRequest{
				SpaceId: spaceId,
				ExportOptions: models.ExportOptions{
					ExportType: exportType,
					TimeZone:   "America/New_York",
					Locale:     "en",
				},
				ShouldExportComments: false,
			},
		},
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(taskRequest); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, constants.NOTION_API_ENQUEUE_URL, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token_v2",
		Value: os.Getenv("NOTION_TOKEN"),
	})
	req.AddCookie(&http.Cookie{
		Name:  "file_token",
		Value: os.Getenv("NOTION_FILE_TOKEN"),
	})

	return req, nil
}

// the first action to get triggered
// adds a task to download the workspace
// sets the task id for the job in the shared data
func (enqueueAction EnqueueAction) Act(s *SharedData) error {
	req, err := enqueueAction.createEnqueueRequest(string(s.ExportType), os.Getenv("NOTION_SPACE_ID"))
	if err != nil {
		return err
	}
	resp, err := enqueueAction.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	var taskResp models.CreateTaskResponseDTO
	if err = json.NewDecoder(resp).Decode(&taskResp); err != nil {
		return err
	}

	if taskResp.TaskId == "" {
		return fmt.Errorf("no task found, enqueue failed")
	}

	// save the taskId to shared data
	s.TaskId = taskResp.TaskId

	return nil
}
