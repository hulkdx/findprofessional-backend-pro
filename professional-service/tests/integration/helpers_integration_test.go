package integration_test

import (
	"encoding/json"
	"net/http/httptest"
)

func Unmarshal(response *httptest.ResponseRecorder, output any) {
	json.Unmarshal(response.Body.Bytes(), output)
}
