package processors

import (
	"fmt"

	"github.com/asishshaji/notion-backup/app/httpclient"
)

type HTMLProcessor struct {
	httpClient *httpclient.HTTPClient
}

func NewHTMLProcessor(client *httpclient.HTTPClient) Processor {
	return &HTMLProcessor{
		httpClient: client,
	}
}

func (hP *HTMLProcessor) Process() error {
	fmt.Println("HTML Processing")
	return nil
}
