package evatr_test

import (
	"testing"

	"github.com/hostwithquantum/go-evatr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExampleClient_ValidateVAT(t *testing.T) {
	client := evatr.NewClient()

	result, err := client.ValidateVAT(t.Context(), requestingVatID, requestedVatID)
	require.NoError(t, err)
	assert.True(t, result.IsValid())

	timestamp, err := result.GetRequestTimestamp()
	require.NoError(t, err)
	assert.False(t, timestamp.IsZero())
	t.Logf("Request timestamp: %s", timestamp)
}

func TestExampleClient_ValidateVATQualified(t *testing.T) {
	client := evatr.NewClient()

	result, err := client.ValidateVATQualified(
		t.Context(),
		requestingVatID,
		requestedVatID,
		exampleCompanyName,
		exampleCity,
		exampleStreet,
		examplePostalCode,
	)
	require.NoError(t, err)
	assert.True(t, result.IsValid())

	t.Log("VAT ID is valid")
	t.Logf("Company name match: %s\n", result.CompanyNameResult)
	t.Logf("City match: %s\n", result.CityResult)
	t.Logf("Street match: %s\n", result.StreetResult)
	t.Logf("Postal code match: %s\n", result.PostalCodeResult)
}

func TestExampleClient_GetEUMemberStates(t *testing.T) {
	client := evatr.NewClient()

	states, err := client.GetEUMemberStates(t.Context())
	require.NoError(t, err)

	t.Log("EU Member States:")
	for _, state := range states {
		status := "available"
		if !state.Available {
			status = "unavailable"
		}
		t.Logf("- %s (%s): %s\n", state.Name, state.Alpha2, status)
	}
}

func TestExampleClient_GetStatusMessages(t *testing.T) {
	client := evatr.NewClient()

	messages, err := client.GetStatusMessages(t.Context())
	require.NoError(t, err)

	t.Log("Status Messages:")
	for _, msg := range messages {
		t.Logf("- %s (HTTP %d): %s\n", msg.Status, msg.HTTPCode, msg.Message)
	}
}

func TestExample_errorHandling(t *testing.T) {
	client := evatr.NewClient()

	_, err := client.ValidateVAT(t.Context(), requestingVatID, "INVALID")
	require.Error(t, err)
	require.True(t, evatr.IsEvatrErr(err))

	evatrErr := err.(*evatr.Error)
	switch evatrErr.Status {
	case evatr.StatusInvalidRequestedVATID:
		t.Logf("The requested VAT ID format is invalid")
	case evatr.StatusVATIDNotAssigned:
		t.Fatal("The VAT ID is not assigned")
	default:
		t.Fatalf("Error: %s\n", evatrErr.Message)
	}
}

func TestExample_verificationResults(t *testing.T) {
	client := evatr.NewClient()

	result, err := client.ValidateVATQualified(
		t.Context(),
		requestingVatID,
		requestedVatID,
		exampleCompanyName,
		exampleCity,
		exampleStreet,
		examplePostalCode,
	)
	require.NoError(t, err)

	checkResult := func(t *testing.T, fieldName, value string) {
		t.Helper()

		switch evatr.VerificationResult(value) {
		case evatr.VerificationMatch:
			t.Logf("%s: matches\n", fieldName)
		case evatr.VerificationMismatch:
			t.Logf("%s: does not match\n", fieldName)
		case evatr.VerificationNotRequested:
			t.Logf("%s: not requested\n", fieldName)
		case evatr.VerificationNotProvided:
			t.Logf("%s: not provided by member state\n", fieldName)
		}
	}

	checkResult(t, "Company name", result.CompanyNameResult)
	checkResult(t, "City", result.CityResult)
	checkResult(t, "Street", result.StreetResult)
	checkResult(t, "Postal code", result.PostalCodeResult)
}

func TestExample_timestampParsing(t *testing.T) {
	client := evatr.NewClient()

	result, err := client.ValidateVAT(t.Context(), requestingVatID, requestedVatID)
	require.NoError(t, err)

	timestamp, err := result.GetRequestTimestamp()
	require.NoError(t, err)
	assert.False(t, timestamp.IsZero())
	t.Logf("Request timestamp: %s", timestamp.Format("2006-01-02 15:04:05 MST"))

	validFrom, err := result.GetValidFrom()
	require.NoError(t, err)
	if !validFrom.IsZero() {
		t.Logf("Valid from: %s", validFrom.Format("2006-01-02"))
	}

	validUntil, err := result.GetValidUntil()
	require.NoError(t, err)
	if !validUntil.IsZero() {
		t.Logf("Valid until: %s", validUntil.Format("2006-01-02"))
	}
}
