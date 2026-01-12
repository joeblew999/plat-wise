// Package wise provides a Go client for the Wise (formerly TransferWise) API.
// API Reference: https://docs.wise.com/api-reference
package wise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// ProductionBaseURL is the base URL for the Wise production API.
	ProductionBaseURL = "https://api.wise.com"
	// SandboxBaseURL is the base URL for the Wise sandbox API.
	SandboxBaseURL = "https://api.sandbox.transferwise.tech"

	defaultTimeout = 30 * time.Second
)

// Client is the Wise API client.
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client

	// Services
	Profiles      *ProfilesService
	Quotes        *QuotesService
	Recipients    *RecipientsService
	Transfers     *TransfersService
	ExchangeRates *ExchangeRatesService
	Balances      *BalancesService
}

// ClientOption is a function that configures the Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithSandbox configures the client to use the sandbox environment.
func WithSandbox() ClientOption {
	return func(c *Client) {
		c.baseURL = SandboxBaseURL
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient creates a new Wise API client.
func NewClient(apiToken string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:  ProductionBaseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.Profiles = &ProfilesService{client: c}
	c.Quotes = &QuotesService{client: c}
	c.Recipients = &RecipientsService{client: c}
	c.Transfers = &TransfersService{client: c}
	c.ExchangeRates = &ExchangeRatesService{client: c}
	c.Balances = &BalancesService{client: c}

	return c
}

// Request performs an HTTP request to the Wise API.
func (c *Client) Request(ctx context.Context, method, path string, query url.Values, body, result interface{}) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, query url.Values, result interface{}) error {
	return c.Request(ctx, http.MethodGet, path, query, nil, result)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.Request(ctx, http.MethodPost, path, nil, body, result)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	return c.Request(ctx, http.MethodPut, path, nil, body, result)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.Request(ctx, http.MethodDelete, path, nil, nil, result)
}
