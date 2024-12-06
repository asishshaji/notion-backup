package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/asishshaji/notion-backup/models"
)

type EnqueueAction struct {
	HttpClient *http.Client
}

func (enqueueAction EnqueueAction) Act(s *SharedData) error {
	taskRequest := models.TaskRequestDTO{
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

	resp, err := enqueueAction.post(marshalledTaskRequest)
	if err != nil {
		return err
	}
	var taskResp models.TaskResponseDTO

	if err = json.Unmarshal(resp, &taskResp); err != nil {
		return err
	}

	// save the taskId to shared data
	s.TaskId = taskResp.TaskId

	return nil
}

func (eA EnqueueAction) post(body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", models.NOTION_API_ENQUEUE_URL, bytes.NewReader(body))
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

	res, err := eA.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing enqueued task id : %s", err)
	}

	return respBody, nil
}
