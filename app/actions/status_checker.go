package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/asishshaji/notion-backup/models"
)

type StatusCheckerAction struct {
	HttpClient *http.Client
}

func (sca StatusCheckerAction) Act(s *SharedData) error {
	// poll the status of the task
	ticker := time.NewTicker(time.Second * 5) // TODO make it configurable

	for {
		select {
		case <-ticker.C:
			status, exportURL, err := sca.getTaskStatus(s.TaskId)
			if err != nil {
				return err
			}
			if status == "success" {
				s.ExportURL = exportURL
			} else if status == "in_progress" {
				fmt.Printf("%s polled, status %s\n", s.TaskId, status)
			}
		case <-time.After(time.Second * 60):
			// something went wrong their side.
			return fmt.Errorf("timed out for task :%s", s.TaskId)
		}
	}
}

func (sca StatusCheckerAction) getTaskStatus(taskId string) (string, string, error) {
	var getTasksDTO models.GetTasksDTO
	var status, exportURL string
	taskIdsReq := struct {
		TaskIds []string `json:"taskIds"`
	}{
		TaskIds: []string{taskId},
	}

	b, err := json.Marshal(taskIdsReq)
	if err != nil {
		return status, exportURL, err
	}

	resp, err := sca.post(b)
	if err != nil {
		return status, exportURL, err
	}

	if err := json.Unmarshal(resp, &getTasksDTO); err != nil {
		return status, exportURL, err
	}

	// TODO add error checking
	result := getTasksDTO.Results[0]
	status = result.State
	exportURL = result.Status.ExportURL

	return status, exportURL, nil
}

// TODO use a custom http client and move these to that
func (sca StatusCheckerAction) post(body []byte) ([]byte, error) {
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

	res, err := sca.HttpClient.Do(req)
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
