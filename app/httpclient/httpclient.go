package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

type HTTPClient struct {
	HttpClient        *http.Client
	NOTION_TOKEN      string
	NOTION_FILE_TOKEN string
	NOTION_SPACE_ID   string
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		HttpClient: &http.Client{Transport: &http.Transport{MaxIdleConns: 10, DisableCompression: true}},
	}
}

func (hC *HTTPClient) Post(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token_v2",
		Value: os.Getenv("NOTION_TOKEN"),
	})
	req.AddCookie(&http.Cookie{
		Name:  "file_token",
		Value: os.Getenv("NOTION_FILE_TOKEN"),
	})

	res, err := hC.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing enqueued task id : %s", err)
	}

	return respBody, nil
}

func (hC *HTTPClient) Get(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token_v2",
		Value: os.Getenv("NOTION_TOKEN"),
	})
	req.AddCookie(&http.Cookie{
		Name:  "file_token",
		Value: os.Getenv("NOTION_FILE_TOKEN"),
	})

	res, err := hC.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing enqueued task id : %s", err)
	}

	return respBody, nil
}
