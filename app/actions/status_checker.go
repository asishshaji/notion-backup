package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/constants"
	"github.com/asishshaji/notion-backup/models"
)

type StatusCheckerAction struct {
	HttpClient *httpclient.HTTPClient
}

func (StatusCheckerAction) String() string {
	return "StatusCheckerAction"
}

// create http request for checking status
func (sca StatusCheckerAction) createStatusRequest(taskId string) (*http.Request, error) {
	taskIdsReq := struct {
		TaskIds []string `json:"taskIds"`
	}{
		TaskIds: []string{taskId},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(taskIdsReq); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, constants.NOTION_API_GET_TASKS_URL, buf)
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

// After enqueueing the task, the function should be triggered to check if zip
// is available for download. Function keeps polling notion api to check if the zip is ready.
// This action sets the downloadUrl in sharedData for the DownloadAction to use
func (sca StatusCheckerAction) Act(s *SharedData) error {
	// poll the status of the task
	ticker := time.NewTicker(time.Second * 10) // TODO make it configurable
	var pollCounter int

	for {
		select {
		case <-ticker.C:
			pollCounter += 1
			fmt.Printf("[%d] polling %s\n", pollCounter, s.TaskId)

			// create request
			req, err := sca.createStatusRequest(s.TaskId)
			if err != nil {
				return err
			}

			// send the http request
			status, exportURL, err := sca.getTaskStatus(req)
			if err != nil {
				return err
			}
			if status == "success" {
				if exportURL == "" {
					return fmt.Errorf("failed to get export url")
				}

				// your workspace is now available to download,
				// move to the next action
				s.ExportURL = exportURL
				return nil
			} else if status == "in_progress" {
				fmt.Printf("%s polled, status %s\n", s.TaskId, status)
				break
			} else {
				return fmt.Errorf("invalid status for task :%s", s.TaskId)
			}
		case <-time.After(time.Second * 120):
			// something went wrong their side.
			return fmt.Errorf("timed out for task :%s", s.TaskId)
		}
	}
}

func (sca StatusCheckerAction) getTaskStatus(req *http.Request) (string, string, error) {
	var status, exportURL string

	// send the request
	resp, err := sca.HttpClient.Do(req)
	if err != nil {
		return status, exportURL, err
	}
	defer resp.Close()

	// parse the response
	var getTasksDTO models.GetTasksDTO
	if err := json.NewDecoder(resp).Decode(&getTasksDTO); err != nil {
		return status, exportURL, err
	}

	if len(getTasksDTO.Results) == 0 {
		return status, exportURL, fmt.Errorf("results should be atleast one")
	}

	result := getTasksDTO.Results[0]
	status = result.State
	exportURL = result.Status.ExportURL

	return status, exportURL, nil
}
