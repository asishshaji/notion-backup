package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/asishshaji/notion-backup/app"
	"github.com/asishshaji/notion-backup/app/httpclient"
	"github.com/asishshaji/notion-backup/app/processors"
	"github.com/asishshaji/notion-backup/constants"
	"github.com/asishshaji/notion-backup/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	httpClient := httpclient.NewHTTPClient()

	// create new app
	app := app.NewApp(httpClient)

	// register all the processors
	app.RegisterProcessor(constants.HtmlExportType, processors.NewHTMLProcessor(httpClient))
	app.RegisterProcessor(constants.MardownExportType, processors.NewMDProcessor(httpClient))

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
			go app.StartProcess(procType, wg)
		}
	}

	// wait till all the goroutines are done
	wg.Wait()

}
