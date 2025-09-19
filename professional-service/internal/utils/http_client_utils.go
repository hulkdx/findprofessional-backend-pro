package utils

import (
	"context"
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
	method,
	url,
	body string,
) (string, error) {
	bodyResponse, err := DoHttpRequestAsReader(ctx, httpClient, method, url, body)
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

func DoHttpRequestAsReader(
	ctx context.Context,
	httpClient *http.Client,
	method,
	url,
	body string,
) (io.ReadCloser, error) {
	bodyReader := strings.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
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
