package processors

import (
	"github.com/asishshaji/notion-backup/app/actions"
	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/models"
)

type MDProcessor struct {
	httpClient *httpclient.HTTPClient
}

type MDSharedData struct{}

func NewMDProcessor(client *httpclient.HTTPClient) Processor {
	return &MDProcessor{
		httpClient: client,
	}
}

func (md *MDProcessor) Actions() []actions.Action {
	return []actions.Action{
		&actions.EnqueueAction{HttpClient: md.httpClient},
		&actions.StatusCheckerAction{HttpClient: md.httpClient},
	}
}

func (md *MDProcessor) Process() error {
	var err error

	s := new(actions.SharedData)
	s.ExportType = models.MardownExportType

	// loop over actions and call act
	for _, action := range md.Actions() {
		err = action.Act(s)
		if err != nil {
			break
		}
	}

	return err
}
