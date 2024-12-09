package httpclient

import (
	"crypto/tls"
	"io"
	"net/http"
)

type HTTPClient struct {
	HttpClient        *http.Client
	NOTION_TOKEN      string
	NOTION_FILE_TOKEN string
	NOTION_SPACE_ID   string
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		HttpClient: &http.Client{Transport: &http.Transport{MaxIdleConns: 10, DisableCompression: true, TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		}}},
	}
}

func (hC *HTTPClient) Do(req *http.Request) (io.ReadCloser, error) {
	res, err := hC.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
