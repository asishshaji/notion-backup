package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/asishshaji/notion-backup/app"
	"github.com/asishshaji/notion-backup/app/processors"
	"github.com/asishshaji/notion-backup/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// create new app
	app := app.NewApp(app.NotionTokens{
		TOKEN:      os.Getenv("NOTION_TOKEN"),
		FILE_TOKEN: os.Getenv("NOTION_FILE_TOKEN"),
		SPACE_ID:   os.Getenv("NOTION_SPACE_ID"),
	}, app.NewHttpClient())

	// register all the processors
	app.RegisterProcessor(models.HtmlExportType, processors.NewHTMLProcessor())
	app.RegisterProcessor(models.MardownExportType, processors.NewMDProcessor())

	flagMap := make(map[*bool]models.ExportType)

	// create flags dynamically for the processors
	for procType := range app.Processors {
		flag := flag.Bool(string(procType), false, fmt.Sprintf("export type %s", procType))
		flagMap[flag] = procType
	}

	flag.Parse()

	wg := new(sync.WaitGroup)
	for flagPtr, procType := range flagMap {
		if *flagPtr {
			wg.Add(1)
			go app.ProcessJob(procType, wg)
		}
	}

	// wait till all the goroutines are done
	wg.Wait()

}
