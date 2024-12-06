package processors

import "fmt"

type MDProcessor struct{}

func NewMDProcessor() Processor {
	return &MDProcessor{}
}

func (md *MDProcessor) Process() error {
	fmt.Println("MD Processing")
	return nil
}
