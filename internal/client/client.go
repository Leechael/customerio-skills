package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func New(region string) (*Client, error) {
	token := os.Getenv("CUSTOMERIO_API_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("CUSTOMERIO_API_TOKEN environment variable is not set")
	}

	baseURL := os.Getenv("CIO_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.customer.io"
		if region == "eu" {
			baseURL = "https://api-eu.customer.io"
		}
	}

	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *Client) do(method, path string, query url.Values, body io.Reader) (json.RawMessage, error) {
	u := c.BaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(data))
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg)
	}

	if len(data) == 0 {
		data = []byte("{}")
	}

	return json.RawMessage(data), nil
}

func (c *Client) Get(path string, query url.Values) (json.RawMessage, error) {
	return c.do(http.MethodGet, path, query, nil)
}

func (c *Client) Post(path string, body any) (json.RawMessage, error) {
	b, err := encodeBody(body)
	if err != nil {
		return nil, err
	}
	return c.do(http.MethodPost, path, nil, b)
}

func (c *Client) Put(path string, body any) (json.RawMessage, error) {
	b, err := encodeBody(body)
	if err != nil {
		return nil, err
	}
	return c.do(http.MethodPut, path, nil, b)
}

func (c *Client) Delete(path string, query url.Values) (json.RawMessage, error) {
	return c.do(http.MethodDelete, path, query, nil)
}

func encodeBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	switch v := body.(type) {
	case json.RawMessage:
		return bytes.NewReader(v), nil
	case []byte:
		return bytes.NewReader(v), nil
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(data), nil
	}
}
