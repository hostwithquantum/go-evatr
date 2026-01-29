package evatr_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hostwithquantum/go-evatr"
	"github.com/stretchr/testify/require"
)

func TestDebugTransport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(evatr.ValidationResponse{
			RequestTimestamp: "2024-01-01T12:00:00Z",
			Status:           evatr.StatusValid,
		}))
	}))
	defer server.Close()

	debugTransport := evatr.NewDebugTransport(nil)
	httpClient := &http.Client{
		Transport: debugTransport,
	}

	client := evatr.NewClient(
		evatr.WithBaseURL(server.URL),
		evatr.WithHTTPClient(httpClient),
	)

	_, err := client.ValidateVAT(t.Context(), "DE123456789", "ATU12345678")
	require.NoError(t, err)
}
