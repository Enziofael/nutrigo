package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPClient interface {
	GET(ctx context.Context, path string, result interface{}) error
	POST(ctx context.Context, path string, body interface{}, result interface{}) error
	PUT(ctx context.Context, path string, body interface{}, result interface{}) error
	PATCH(ctx context.Context, path string, body interface{}, result interface{}) error
	DELETE(ctx context.Context, path string, result interface{}) error
	HEAD(ctx context.Context, path string) error
	DO(ctx context.Context, method, path string, body interface{}, result interface{}) error
}

type Client struct {
	BaseUrl string
	Client  *http.Client
	headers map[string]string
}

func NewClient(baseUrl string, cfg http.HTTP2Config, apikey string) *Client {
	return &Client{}
}

func (c *Client) GET(ctx context.Context, path string, result interface{}) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) POST(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

func (c *Client) PUT(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.do(ctx, http.MethodPut, path, body, result)
}

func (c *Client) PATCH(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.do(ctx, http.MethodPatch, path, body, result)
}

func (c *Client) DELETE(ctx context.Context, path string, result interface{}) error {
	return c.do(ctx, http.MethodDelete, path, nil, result)
}

func (c *Client) HEAD(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodHead, path, nil, nil)
}

func (c *Client) DO(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	return c.do(ctx, method, path, body, result)
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {

	url := c.buildURL(path)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marchal request body: %v", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %v", err)
	}

	req.Header.Set("Content-Type", c.headers["Content-Type"])
	req.Header.Set("Accept", c.headers["Accept"])
	req.Header.Set("User-Agent", c.headers["User-Agent"])

	if c.headers["X-API-Key"] != "" {
		req.Header.Set("X-API-Key", c.headers["X-API-Key"])
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Читаем тело ошибки для диагностики
		bodyBytes, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(bodyBytes),
			Err:        fmt.Errorf("unexpected status code"),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) buildURL(path string) string {
	path = strings.Trim(path, "/")
	return fmt.Sprintf("%v/%v", c.BaseUrl, path)
}

func (c *Client) Close() error {
	c.Client.CloseIdleConnections()
	return nil
}
