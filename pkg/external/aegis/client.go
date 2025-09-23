package aegis

import (
	pkgHttp "cef/pkg/http"
	"context"
	"fmt"
	"net/http"
	"strings"
)

var defaultClient *Client

func DefaultClient() *Client {
	return defaultClient
}

func SetDefault(cli *Client) {
	defaultClient = cli
}

type Client struct {
	addr string
}

func NewAegisClient(addr string) *Client {
	return &Client{strings.TrimSuffix(addr, "/")}
}

// GetConfig /config/{namespace}/{service}/get?keys=xxx
func (c *Client) GetConfig(serviceName string, keys ...string) (map[string]any, error) {
	path := fmt.Sprintf("/config/cef/%s/get", serviceName)
	if len(keys) > 0 {
		path += "?"
		for _, key := range keys {
			path += "keys=" + key + "&"
		}
		strings.TrimSuffix(path, "&")
	}
	var resp ConfigGetResponse
	if err := pkgHttp.DoWithJsonResult(context.Background(), pkgHttp.MustNewRequest(http.MethodGet, c.addr+path, nil), &resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("result code is not 0, msg: %s", resp.Msg)
	}
	return resp.Data.ConfigMap, nil
}

func (c *Client) GetConfigWithResult(serviceName string, result any, keys ...string) error {
	path := fmt.Sprintf("/config/cef/%s/get", serviceName)
	if len(keys) > 0 {
		path += "?"
		for _, key := range keys {
			path += "keys=" + key + "&"
		}
		path = strings.TrimSuffix(path, "&")
	}
	return pkgHttp.DoWithJsonResult(context.Background(), pkgHttp.MustNewRequest(http.MethodGet, c.addr+path, nil), result)
}
