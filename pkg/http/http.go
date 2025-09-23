package pkgHttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func MustNewRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(fmt.Errorf("http.NewRequest(method, url, body) failed, %w", err))
	}
	return req
}

func DoWithJsonResult(ctx context.Context, req *http.Request, result interface{}) error {
	if req == nil || result == nil {
		return nil
	}
	req.WithContext(ctx)
	if (req.Method == http.MethodPost || req.Method == http.MethodPut) && req.Header.Get(HeaderContentType) == "" {
		req.Header.Set(HeaderContentType, "application/json;charset=utf-8")
	}
	fmt.Printf("http.%s(req): %s\n", req.Method, req.URL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do(req) failed, %w", err)
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	if err = decoder.Decode(result); err != nil {
		return fmt.Errorf("decoder.Decode(result) failed, %w", err)
	}
	return nil
}
