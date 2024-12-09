package actions

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/asishshaji/notion-backup/app/httpclient"
)

type DownloaderAction struct {
	HttpClient *httpclient.HTTPClient
}

func (DownloaderAction) String() string {
	return "DownloaderAction"
}

func (dA DownloaderAction) createDownloadRequest(exportUrl string) (*http.Request, error) {
	req, _ := http.NewRequest(http.MethodGet, exportUrl, nil)
	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token_v2",
		Value: os.Getenv("NOTION_TOKEN"),
	})
	req.AddCookie(&http.Cookie{
		Name:  "file_token",
		Value: os.Getenv("NOTION_FILE_TOKEN"),
	})
	return req, nil
}

// called after checking the status of the task enqueued to notion
// downloads the zip to tmp directory and extracts it to current directory
func (dA DownloaderAction) Act(s *SharedData) error {
	fmt.Printf("downloading from %s -> export type: %s", s.ExportURL, s.ExportType)
	req, err := dA.createDownloadRequest(s.ExportURL)
	if err != nil {
		return err
	}
	// download the zip
	resp, err := dA.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	// create the zip file
	fileName := fmt.Sprintf("/tmp/%s.zip", s.ExportType)
	outFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer outFile.Close()
	// copy contents to zip file from the response
	_, err = io.Copy(outFile, resp)
	if err != nil {
		return err
	}

	// set the downloaded file path for the next actions to use
	s.DownloadedFilePath = fileName
	fmt.Printf("downloaded %s to %s\n", s.ExportType, fileName)

	return nil
}
