package processors

import (
	"fmt"

	"github.com/asishshaji/notion-backup/app/actions"
	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/constants"
)

type MDProcessor struct {
	httpClient *httpclient.HTTPClient
	actions    []actions.Action
}

func NewMDProcessor(client *httpclient.HTTPClient) Processor {
	return &MDProcessor{
		httpClient: client,
	}
}

func (md *MDProcessor) initialiseActions() {
	// define the sequence of actions the process should follow
	md.actions = []actions.Action{
		&actions.EnqueueAction{HttpClient: md.httpClient},
		&actions.StatusCheckerAction{HttpClient: md.httpClient},
		&actions.DownloaderAction{HttpClient: md.httpClient},
		&actions.ExtractorAction{},
	}
}

func (md *MDProcessor) Actions() []actions.Action {
	return md.actions
}

func (md *MDProcessor) Process() error {
	var err error
	md.initialiseActions()

	s := new(actions.SharedData)
	s.ExportType = constants.MardownExportType

	// loop over actions and call act
	for i, action := range md.Actions() {
		fmt.Printf("Executing action:%d %s\n", i+1, action)
		err = action.Act(s)
		if err != nil {
			break
		}
	}

	return err
}
