package pkgHttp

import (
	"context"
	"net/http"
	"testing"
)

func TestDoWithJsonResult(t *testing.T) {
	result := make(map[string]any)
	if err := DoWithJsonResult(context.Background(), MustNewRequest(http.MethodGet, `http://0.0.0.0:11470/config/cef/browser-config/get?keys=["test"]`, nil), &result); err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
