package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func DoHttpRequest(
	ctx context.Context,
	httpClient *http.Client,
	method,
	url,
	body string,
) (string, error) {
	bodyReader := strings.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status %d from %s", res.StatusCode, url)
	}
	byt, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}
