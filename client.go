package evatr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	// Default API endpoint for the eVatR service
	DefaultBaseURL = "https://api.evatr.vies.bzst.de/app"

	// Default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

var version = func() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "devel"
	}
	return info.Main.Version
}()

// Client is the eVATR API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithBaseURL sets the base URL for the API.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the timeout for HTTP requests.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// NewClient returns a new eVATR API client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request and handles common error responses.
func (c *Client) doRequest(ctx context.Context, method, path string, body any, result any) error {
	var reqBody io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		reqBody = &buf
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("go-evatr/%s (+https://github.com/hostwithquantum/go-evatr)", version))
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}
		return nil
	}

	return c.handleErrorResponse(resp.StatusCode, resp.Body)
}

// handleErrorResponse converts HTTP error responses into typed errors.
func (c *Client) handleErrorResponse(statusCode int, body io.Reader) error {
	var errResp ErrorResponse
	var status, message string

	if err := json.NewDecoder(body).Decode(&errResp); err == nil {
		status = errResp.Status
		message = errResp.Message
	}

	if message == "" {
		message = getDefaultErrorMessage(statusCode, status)
	}

	return &Error{
		StatusCode: statusCode,
		Status:     status,
		Message:    message,
	}
}

// getDefaultErrorMessage returns a default error message based on status code and evatr status.
func getDefaultErrorMessage(statusCode int, status string) string {
	switch statusCode {
	case 400:
		return "Bad request: Invalid input parameters"
	case 403:
		return "Forbidden: Not authorized to perform this request"
	case 404:
		return "Not found: VAT ID not found or requesting VAT ID invalid"
	case 500:
		return "Internal server error: Processing temporarily not possible"
	case 503:
		return "Service unavailable: Please try again later"
	default:
		return fmt.Sprintf("Unexpected error (HTTP %d): %s", statusCode, status)
	}
}
