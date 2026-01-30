package evatr_test

import (
	"testing"

	"github.com/hostwithquantum/go-evatr"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

// implement a health check
func (s *IntegrationSuite) SetupSuite() {
	client := evatr.NewClient()
	_, err := client.GetStatusMessages(s.T().Context())
	if err == nil {
		return
	}

	if !evatr.IsEvatrErr(err) {
		s.T().Fatalf("got an error, marking this as failed: %+v", err)
	}

	eVatRErr := err.(*evatr.Error)
	if eVatRErr.StatusCode == 503 {
		s.T().Skip("eVatR API has maintenance")
	}
	s.T().Fatalf("unexpected err (eVatR): %+v", eVatRErr)
}

func (s *IntegrationSuite) TestValidateVAT() {
	client := evatr.NewClient()

	result, err := client.ValidateVAT(s.T().Context(), requestingVatID, requestedVatID)
	s.Require().NoError(err)
	s.Assert().True(result.IsValid())

	timestamp, err := result.GetRequestTimestamp()
	s.Require().NoError(err)
	s.Assert().False(timestamp.IsZero())
	s.T().Logf("Request timestamp: %s", timestamp)
}

func (s *IntegrationSuite) TestValidateVATQualified() {
	client := evatr.NewClient()

	result, err := client.ValidateVATQualified(
		s.T().Context(),
		requestingVatID,
		requestedVatID,
		exampleCompanyName,
		exampleCity,
		exampleStreet,
		examplePostalCode,
	)
	s.Require().NoError(err)
	s.Assert().True(result.IsValid())

	s.T().Log("VAT ID is valid")
	s.T().Logf("Company name match: %s\n", result.CompanyNameResult)
	s.T().Logf("City match: %s\n", result.CityResult)
	s.T().Logf("Street match: %s\n", result.StreetResult)
	s.T().Logf("Postal code match: %s\n", result.PostalCodeResult)
}

func (s *IntegrationSuite) TestGetEUMemberStates() {
	client := evatr.NewClient()

	states, err := client.GetEUMemberStates(s.T().Context())
	s.Require().NoError(err)

	s.T().Log("EU Member States:")
	for _, state := range states {
		status := "available"
		if !state.Available {
			status = "unavailable"
		}
		s.T().Logf("- %s (%s): %s\n", state.Name, state.Alpha2, status)
	}
}

func (s *IntegrationSuite) TestGetStatusMessages() {
	client := evatr.NewClient()

	messages, err := client.GetStatusMessages(s.T().Context())
	s.Require().NoError(err)

	s.T().Log("Status Messages:")
	for _, msg := range messages {
		s.T().Logf("- %s (HTTP %d): %s\n", msg.Status, msg.HTTPCode, msg.Message)
	}
}

func (s *IntegrationSuite) TestErrorHandling() {
	client := evatr.NewClient()

	_, err := client.ValidateVAT(s.T().Context(), requestingVatID, "INVALID")
	s.Require().Error(err)
	s.Require().True(evatr.IsEvatrErr(err))

	evatrErr := err.(*evatr.Error)
	switch evatrErr.Status {
	case evatr.StatusInvalidRequestedVATID:
		s.T().Logf("The requested VAT ID format is invalid")
	case evatr.StatusVATIDNotAssigned:
		s.T().Fatal("The VAT ID is not assigned")
	default:
		s.T().Fatalf("Error: %s\n", evatrErr.Message)
	}
}

func (s *IntegrationSuite) TestVerificationResults() {
	client := evatr.NewClient()

	result, err := client.ValidateVATQualified(
		s.T().Context(),
		requestingVatID,
		requestedVatID,
		exampleCompanyName,
		exampleCity,
		exampleStreet,
		examplePostalCode,
	)
	s.Require().NoError(err)

	checkResult := func(fieldName string, value evatr.VerificationResult) {
		s.T().Helper()

		switch evatr.VerificationResult(value) {
		case evatr.VerificationMatch:
			s.T().Logf("%s: matches\n", fieldName)
		case evatr.VerificationMismatch:
			s.T().Logf("%s: does not match\n", fieldName)
		case evatr.VerificationNotRequested:
			s.T().Logf("%s: not requested\n", fieldName)
		case evatr.VerificationNotProvided:
			s.T().Logf("%s: not provided by member state\n", fieldName)
		}
	}

	checkResult("Company name", result.CompanyNameResult)
	checkResult("City", result.CityResult)
	checkResult("Street", result.StreetResult)
	checkResult("Postal code", result.PostalCodeResult)
}

func (s *IntegrationSuite) TestTimestampParsing() {
	client := evatr.NewClient()

	result, err := client.ValidateVAT(s.T().Context(), requestingVatID, requestedVatID)
	s.Require().NoError(err)

	timestamp, err := result.GetRequestTimestamp()
	s.Require().NoError(err)
	s.Assert().False(timestamp.IsZero())
	s.T().Logf("Request timestamp: %s", timestamp.Format("2006-01-02 15:04:05 MST"))

	validFrom, err := result.GetValidFrom()
	s.Require().NoError(err)
	if !validFrom.IsZero() {
		s.T().Logf("Valid from: %s", validFrom.Format("2006-01-02"))
	}

	validUntil, err := result.GetValidUntil()
	s.Require().NoError(err)
	if !validUntil.IsZero() {
		s.T().Logf("Valid until: %s", validUntil.Format("2006-01-02"))
	}
}
