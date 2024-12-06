package app

import (
	"fmt"
	"sync"

	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/app/processors"
	"github.com/asishshaji/notion-backup/models"
)

type App struct {
	Client     *httpclient.HTTPClient
	Processors map[models.ExportType]processors.Processor
}

func NewApp(httpClient *httpclient.HTTPClient) *App {
	return &App{
		Client:     httpClient,
		Processors: make(map[models.ExportType]processors.Processor),
	}
}

func (app *App) StartProcess(exportType models.ExportType, wg *sync.WaitGroup) {
	defer wg.Done()
	// get the processor for the export type
	proc := app.Processors[exportType]

	// call the process function for the processor
	if err := proc.Process(); err != nil {
		fmt.Printf("error processing for %s: %s", exportType, err)
		return
	}
}

func (app *App) RegisterProcessor(exportType models.ExportType, processor processors.Processor) {
	app.Processors[exportType] = processor
}
