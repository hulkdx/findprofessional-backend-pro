package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func CreateDefaultAppHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func DoHttpRequestAsString(
	ctx context.Context,
	httpClient *http.Client,
	method string,
	url string,
	body string,
) (string, error) {
	bodyReader := strings.NewReader(body)
	bodyResponse, err := DoHttpRequestAsReader(ctx, httpClient, method, url, bodyReader)
	if err != nil {
		return "", err
	}
	defer bodyResponse.Close()
	byt, err := io.ReadAll(bodyResponse)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

func DoHttpRequestAsStruct(
	ctx context.Context,
	httpClient *http.Client,
	method string,
	url string,
	body any,
	responseStruct any,
) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}
	bodyResponse, err := DoHttpRequestAsReader(ctx, httpClient, method, url, &buf)
	if err != nil {
		return err
	}
	defer bodyResponse.Close()
	return json.NewDecoder(bodyResponse).Decode(responseStruct)
}

func DoHttpRequestAsReader(
	ctx context.Context,
	httpClient *http.Client,
	method string,
	url string,
	body io.Reader,
) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d from %s", res.StatusCode, url)
	}
	return res.Body, nil
}
