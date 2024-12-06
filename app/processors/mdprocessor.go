package processors

import (
	"github.com/asishshaji/notion-backup/app/actions"
	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/constants"
)

type MDProcessor struct {
	httpClient *httpclient.HTTPClient
	actions    []actions.Action
}

type MDSharedData struct{}

func NewMDProcessor(client *httpclient.HTTPClient) Processor {
	return &MDProcessor{
		httpClient: client,
	}
}

func (md *MDProcessor) SetActions() {
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

	s := new(actions.SharedData)
	s.ExportType = constants.MardownExportType

	// loop over actions and call act
	for _, action := range md.Actions() {
		err = action.Act(s)
		if err != nil {
			break
		}
	}

	return err
}
