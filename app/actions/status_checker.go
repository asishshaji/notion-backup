package actions

import "net/http"

type StatusCheckerAction struct {
	HttpClient *http.Client
}

func (sca StatusCheckerAction) Act(s *SharedData) error {
	return nil
}

func (sca StatusCheckerAction) getTaskStatus(taskId string) (string, error) {
	return "", nil
}
