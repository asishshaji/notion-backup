package processors

import "fmt"

type HTMLProcessor struct{}

func NewHTMLProcessor() Processor {
	return &HTMLProcessor{}
}

func (hP *HTMLProcessor) Process() error {
	fmt.Println("HTML Processing")
	return nil
}
