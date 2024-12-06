package app

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/asishshaji/notion-backup/app/processors"
	"github.com/asishshaji/notion-backup/models"
)

type NotionTokens struct {
	TOKEN      string
	FILE_TOKEN string
	SPACE_ID   string
}

func NewHttpClient() *http.Client {
	return &http.Client{Transport: &http.Transport{MaxIdleConns: 10, DisableCompression: true}}
}

type App struct {
	Tokens     NotionTokens
	Client     *http.Client
	Processors map[models.ExportType]processors.Processor
}

func NewApp(tokens NotionTokens, httpClient *http.Client) *App {
	return &App{
		Tokens:     tokens,
		Client:     httpClient,
		Processors: make(map[models.ExportType]processors.Processor),
	}
}

func (app *App) ProcessJob(exportType models.ExportType, wg *sync.WaitGroup) {
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
