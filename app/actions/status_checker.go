package actions

import (
	"encoding/json"
	"fmt"
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

func (sca StatusCheckerAction) Act(s *SharedData) error {
	// poll the status of the task
	ticker := time.NewTicker(time.Second * 20) // TODO make it configurable
	var pollCounter int
	for {
		select {
		case <-ticker.C:
			pollCounter += 1
			fmt.Printf("[%d] polling %s\n", pollCounter, s.TaskId)
			status, exportURL, err := sca.getTaskStatus(s.TaskId)
			if err != nil {
				return err
			}
			if status == "success" {
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

	resp, err := sca.HttpClient.Post(constants.NOTION_API_GET_TASKS_URL, b)
	if err != nil {
		return status, exportURL, err
	}

	if err := json.Unmarshal(resp, &getTasksDTO); err != nil {
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
