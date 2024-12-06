package actions

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/asishshaji/notion-backup/app/httpclient"
)

type DownloaderAction struct {
	HttpClient *httpclient.HTTPClient
}

func (dA DownloaderAction) Act(s *SharedData) error {
	fmt.Printf("downloading from %s -> export type: %s", s.ExportURL, s.ExportType)
	resp, err := dA.HttpClient.Get(s.ExportURL)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("/tmp/%s.zip", s.ExportType)
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, bytes.NewReader(resp))
	if err != nil {
		return err
	}

	s.DownloadedFilePath = fileName
	fmt.Printf("downloaded %s to %s", s.ExportType, fileName)

	return nil
}
