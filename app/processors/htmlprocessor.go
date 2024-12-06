package processors

import (
	"github.com/asishshaji/notion-backup/app/actions"
	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/models"
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

func (md *HTMLProcessor) SetActions() {
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

	s := new(actions.SharedData)
	s.ExportType = models.HtmlExportType

	// loop over actions and call act
	for _, action := range hP.Actions() {
		err = action.Act(s)
		if err != nil {
			break
		}
	}

	return err
}
