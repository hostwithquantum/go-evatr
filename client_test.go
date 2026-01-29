package evatr_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hostwithquantum/go-evatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient verifies client initialization with various options
func TestNewClient(t *testing.T) {
	t.Run("default client", func(t *testing.T) {
		client := evatr.NewClient()
		assert.NotNil(t, client)
	})

	t.Run("custom base URL", func(t *testing.T) {
		customURL := "https://test.example.com"
		client := evatr.NewClient(evatr.WithBaseURL(customURL))
		assert.NotNil(t, client)
	})

	t.Run("custom timeout", func(t *testing.T) {
		customTimeout := 60 * time.Second
		client := evatr.NewClient(evatr.WithTimeout(customTimeout))
		assert.NotNil(t, client)
	})

	t.Run("custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{Timeout: 5 * time.Second}
		client := evatr.NewClient(evatr.WithHTTPClient(customClient))
		assert.NotNil(t, client)
	})
}

// TestValidateVAT tests simple VAT ID validation
func TestValidateVAT(t *testing.T) {
	t.Run("successful validation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/abfrage", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var req evatr.ValidationRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			require.NoError(t, err)

			assert.Equal(t, "DE123456789", req.RequestingVATID)
			assert.Equal(t, "ATU12345678", req.RequestedVATID)

			resp := evatr.ValidationResponse{
				RequestTimestamp: time.Now().Format(time.RFC3339),
				Status:           evatr.StatusValid,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := evatr.NewClient(evatr.WithBaseURL(server.URL))
		result, err := client.ValidateVAT(t.Context(), "DE123456789", "ATU12345678")
		require.NoError(t, err)

		assert.Equal(t, evatr.StatusValid, result.Status)
		assert.True(t, result.IsValid())
	})

	t.Run("missing requesting VAT ID", func(t *testing.T) {
		client := evatr.NewClient()
		_, err := client.ValidateVAT(t.Context(), "", "ATU12345678")
		require.Error(t, err)
	})

	t.Run("missing requested VAT ID", func(t *testing.T) {
		client := evatr.NewClient()
		_, err := client.ValidateVAT(t.Context(), "DE123456789", "")
		require.Error(t, err)
	})

	t.Run("VAT ID not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(evatr.ErrorResponse{
				Status:  evatr.StatusVATIDNotAssigned,
				Message: "Die angefragte USt-IdNr. ist zum Anfragezeitpunkt nicht vergeben.",
			})
		}))
		defer server.Close()

		client := evatr.NewClient(evatr.WithBaseURL(server.URL))
		_, err := client.ValidateVAT(t.Context(), "DE123456789", "ATU99999999")
		require.Error(t, err)

		assert.True(t, evatr.IsEvatrErr(err))
		evatrErr := err.(*evatr.Error)
		assert.Equal(t, 404, evatrErr.StatusCode)
		assert.Equal(t, evatr.StatusVATIDNotAssigned, evatrErr.Status)
	})
}

// TestValidateVATQualified tests qualified VAT ID validation
func TestValidateVATQualified(t *testing.T) {
	t.Run("successful qualified validation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req evatr.ValidationRequest
			json.NewDecoder(r.Body).Decode(&req)

			assert.NotEmpty(t, req.CompanyName)
			assert.NotEmpty(t, req.City)

			resp := evatr.ValidationResponse{
				RequestTimestamp:  time.Now().Format(time.RFC3339),
				Status:            evatr.StatusValid,
				CompanyNameResult: string(evatr.VerificationMatch),
				CityResult:        string(evatr.VerificationMatch),
				StreetResult:      string(evatr.VerificationMatch),
				PostalCodeResult:  string(evatr.VerificationMatch),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := evatr.NewClient(evatr.WithBaseURL(server.URL))
		result, err := client.ValidateVATQualified(
			t.Context(),
			"DE123456789",
			"ATU12345678",
			"Musterhaus GmbH & Co KG",
			"Musterort",
			"Musterstrasse 22",
			"12345",
		)
		require.NoError(t, err)

		assert.Equal(t, string(evatr.VerificationMatch), result.CompanyNameResult)
		assert.Equal(t, string(evatr.VerificationMatch), result.CityResult)
	})

	t.Run("missing company name", func(t *testing.T) {
		client := evatr.NewClient()
		_, err := client.ValidateVATQualified(t.Context(), "DE123456789", "ATU12345678", "", "Berlin", "", "")
		require.Error(t, err)
	})

	t.Run("missing city", func(t *testing.T) {
		client := evatr.NewClient()
		_, err := client.ValidateVATQualified(t.Context(), "DE123456789", "ATU12345678", "Test GmbH", "", "", "")
		require.Error(t, err)
	})

	t.Run("data mismatch", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := evatr.ValidationResponse{
				RequestTimestamp:  time.Now().Format(time.RFC3339),
				Status:            evatr.StatusValid,
				CompanyNameResult: string(evatr.VerificationMismatch),
				CityResult:        string(evatr.VerificationMatch),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		client := evatr.NewClient(evatr.WithBaseURL(server.URL))
		result, err := client.ValidateVATQualified(
			t.Context(),
			"DE123456789",
			"ATU12345678",
			"Wrong Company Name",
			"Musterort",
			"",
			"",
		)
		require.NoError(t, err)

		assert.Equal(t, string(evatr.VerificationMismatch), result.CompanyNameResult)
	})
}

// TestGetStatusMessages tests retrieving status messages
func TestGetStatusMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/info/statusmeldungen", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		messages := []evatr.StatusMessage{
			{
				Status:   evatr.StatusValid,
				Category: "Erfolg",
				HTTPCode: 200,
				Message:  "Die angefragte Ust-IdNr. ist zum Anfragezeitpunkt gültig.",
			},
			{
				Status:   evatr.StatusVATIDNotAssigned,
				Category: "Fehler",
				HTTPCode: 404,
				Message:  "Die angefragte USt-IdNr. ist zum Anfragezeitpunkt nicht vergeben.",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}))
	defer server.Close()

	client := evatr.NewClient(evatr.WithBaseURL(server.URL))
	messages, err := client.GetStatusMessages(t.Context())
	require.NoError(t, err)

	assert.Equal(t, 2, len(messages))
	assert.Equal(t, evatr.StatusValid, messages[0].Status)
}

// TestGetEUMemberStates tests retrieving EU member states
func TestGetEUMemberStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/info/eu_mitgliedstaaten", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		states := []evatr.EUMemberState{
			{Alpha2: "DE", Name: "Deutschland", Available: true},
			{Alpha2: "AT", Name: "Österreich", Available: true},
			{Alpha2: "FR", Name: "Frankreich", Available: false},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(states)
	}))
	defer server.Close()

	client := evatr.NewClient(evatr.WithBaseURL(server.URL))
	states, err := client.GetEUMemberStates(t.Context())
	require.NoError(t, err)

	assert.Equal(t, 3, len(states))
	assert.Equal(t, "DE", states[0].Alpha2)
	assert.True(t, states[0].Available)
	assert.False(t, states[2].Available)
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name           string
		statusCode     int
		responseBody   evatr.ErrorResponse
		expectedStatus string
		expectedCode   int
	}{
		{
			name:       "bad request - missing field",
			statusCode: 400,
			responseBody: evatr.ErrorResponse{
				Status:  evatr.StatusMissingRequiredField,
				Message: "Mindestens eins der Pflichtfelder ist nicht besetzt.",
			},
			expectedStatus: evatr.StatusMissingRequiredField,
			expectedCode:   400,
		},
		{
			name:       "forbidden - not authorized",
			statusCode: 403,
			responseBody: evatr.ErrorResponse{
				Status:  evatr.StatusNotAuthorizedDE,
				Message: "Die anfragende DE USt-IdNr. ist nicht berechtigt.",
			},
			expectedStatus: evatr.StatusNotAuthorizedDE,
			expectedCode:   403,
		},
		{
			name:       "service unavailable",
			statusCode: 503,
			responseBody: evatr.ErrorResponse{
				Status:  evatr.StatusServiceUnavailable1,
				Message: "Eine Bearbeitung Ihrer Anfrage ist zurzeit nicht möglich.",
			},
			expectedStatus: evatr.StatusServiceUnavailable1,
			expectedCode:   503,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				require.NoError(t, json.NewEncoder(w).Encode(tc.responseBody))
			}))
			defer server.Close()

			client := evatr.NewClient(evatr.WithBaseURL(server.URL))
			_, err := client.ValidateVAT(t.Context(), "DE123456789", "ATU12345678")
			require.Error(t, err)

			assert.True(t, evatr.IsEvatrErr(err))
			evatrErr := err.(*evatr.Error)
			assert.Equal(t, tc.expectedCode, evatrErr.StatusCode)
			assert.Equal(t, tc.expectedStatus, evatrErr.Status)
		})
	}
}

// TestValidationResponseMethods tests helper methods on ValidationResponse
func TestValidationResponseMethods(t *testing.T) {
	t.Run("IsValid", func(t *testing.T) {
		tests := []struct {
			status   string
			expected bool
		}{
			{evatr.StatusValid, true},
			{evatr.StatusValidWithSpecialCase, true},
			{evatr.StatusNotYetValid, false},
			{evatr.StatusNoLongerValid, false},
			{evatr.StatusVATIDNotAssigned, false},
		}

		for _, tt := range tests {
			resp := &evatr.ValidationResponse{Status: tt.status}
			assert.Equal(t, tt.expected, resp.IsValid(), "status %s", tt.status)
		}
	})

	t.Run("GetRequestTimestamp", func(t *testing.T) {
		now := time.Now()
		resp := &evatr.ValidationResponse{
			RequestTimestamp: now.Format(time.RFC3339),
		}
		parsed, err := resp.GetRequestTimestamp()
		require.NoError(t, err)
		assert.Equal(t, now.Unix(), parsed.Unix())
	})

	t.Run("GetValidFrom with value", func(t *testing.T) {
		date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		resp := &evatr.ValidationResponse{
			ValidFrom: date.Format(time.RFC3339),
		}
		parsed, err := resp.GetValidFrom()
		require.NoError(t, err)
		assert.Equal(t, date.Unix(), parsed.Unix())
	})

	t.Run("GetValidFrom empty", func(t *testing.T) {
		resp := &evatr.ValidationResponse{}
		parsed, err := resp.GetValidFrom()
		require.NoError(t, err)
		assert.True(t, parsed.IsZero())
	})
}
