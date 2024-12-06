package processors

import (
	"fmt"

	"github.com/asishshaji/notion-backup/app/actions"
	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/constants"
)

type HTMLProcessor struct {
	httpClient *httpclient.HTTPClient
	actions    []actions.Action
}

func NewHTMLProcessor(client *httpclient.HTTPClient) Processor {
	return &HTMLProcessor{
		httpClient: client,
	}
}

func (md *HTMLProcessor) initialiseActions() {
	// define the sequence of actions the process should follow
	md.actions = []actions.Action{
		&actions.EnqueueAction{HttpClient: md.httpClient},
		&actions.StatusCheckerAction{HttpClient: md.httpClient},
		&actions.DownloaderAction{HttpClient: md.httpClient},
		&actions.ExtractorAction{},
	}
}

func (md *HTMLProcessor) Actions() []actions.Action {
	return md.actions
}

func (hP *HTMLProcessor) Process() error {
	var err error
	hP.initialiseActions()

	s := new(actions.SharedData)
	s.ExportType = constants.HtmlExportType

	// loop over actions and call act
	for i, action := range hP.Actions() {
		fmt.Printf("Executing action:%d %s\n", i+1, action)
		err = action.Act(s)
		if err != nil {
			break
		}
	}

	return err
}
