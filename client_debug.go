package evatr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

type debugTransport struct {
	transport http.RoundTripper
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		fmt.Printf("Error dumping request: %v\n", err)
	} else {
		fmt.Printf("=== REQUEST ===\n%s\n", string(reqDump))
	}

	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		fmt.Printf("=== ERROR ===\n%v\n", err)
		return nil, err
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Printf("Error dumping response: %v\n", err)
	} else {
		fmt.Printf("=== RESPONSE ===\n%s\n", string(respDump))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return resp, nil
}

// NewDebugTransport returns a transport that logs all requests and responses.
func NewDebugTransport(transport http.RoundTripper) http.RoundTripper {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &debugTransport{transport: transport}
}
