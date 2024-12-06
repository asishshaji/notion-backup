package processors

import (
	"fmt"
	"net/http"
)

type HTMLProcessor struct {
	httpClient *http.Client
}

func NewHTMLProcessor(client *http.Client) Processor {
	return &HTMLProcessor{
		httpClient: client,
	}
}

func (hP *HTMLProcessor) Process() error {
	fmt.Println("HTML Processing")
	return nil
}
