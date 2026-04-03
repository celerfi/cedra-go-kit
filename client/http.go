package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type APIError struct {
	Message     string `json:"message"`
	ErrorCode   string `json:"error_code"`
	VmErrorCode *int   `json:"vm_error_code,omitempty"`
	StatusCode  int    `json:"-"`
}

func (e *APIError) Error() string {
	if e.VmErrorCode != nil {
		return fmt.Sprintf("cedra API error [%d] %s: %s (vm_error_code=%d)", e.StatusCode, e.ErrorCode, e.Message, *e.VmErrorCode)
	}
	return fmt.Sprintf("cedra API error [%d] %s: %s", e.StatusCode, e.ErrorCode, e.Message)
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *Client) NodeURL() string    { return c.cfg.NodeURL }
func (c *Client) FaucetURL() string  { return c.cfg.FaucetURL }
func (c *Client) Network() Network   { return c.cfg.Network }

func (c *Client) Get(ctx context.Context, path string, params url.Values, result interface{}) error {
	u := c.cfg.NodeURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	return c.do(req, result)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.NodeURL+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.do(req, result)
}

func (c *Client) PostBCS(ctx context.Context, path string, body []byte, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.NodeURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x.cedra.signed_transaction+bcs")
	req.Header.Set("Accept", "application/json")
	return c.do(req, result)
}

func (c *Client) PostFaucet(ctx context.Context, path string, body interface{}, result interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.FaucetURL+path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.do(req, result)
}

func (c *Client) do(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		apiErr.StatusCode = resp.StatusCode
		if jsonErr := json.Unmarshal(body, &apiErr); jsonErr != nil {
			apiErr.Message = string(body)
			apiErr.ErrorCode = "unknown"
		}
		return &apiErr
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}
