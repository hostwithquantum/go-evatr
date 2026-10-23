package evatr

import (
	"context"
	"fmt"
	"strings"

	t "github.com/anyascii/go"
)

// ValidateVAT validates a VAT ID without company data verification.
func (c *Client) ValidateVAT(ctx context.Context, requestingVATID, requestedVATID string) (*ValidationResponse, error) {
	if err := validateRequestingVAT(requestingVATID); err != nil {
		return nil, err
	}
	if err := validateRequestedVAT(requestedVATID); err != nil {
		return nil, err
	}

	req := &ValidationRequest{
		RequestingVATID: requestingVATID,
		RequestedVATID:  requestedVATID,
	}

	var resp ValidationResponse
	if err := c.doRequest(ctx, "POST", "/v1/abfrage", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ValidateVATQualified validates a VAT ID with company data verification.
// Compares provided company information with registered data.
func (c *Client) ValidateVATQualified(ctx context.Context, requestingVATID, requestedVATID, companyName, city, street, postalCode string) (*ValidationResponse, error) {
	if err := validateRequestingVAT(requestingVATID); err != err {
		return nil, err
	}
	if err := validateRequestedVAT(requestedVATID); err != nil {
		return nil, err
	}
	if companyName == "" {
		return nil, fmt.Errorf("company name is required for qualified validation")
	}
	if city == "" {
		return nil, fmt.Errorf("city is required for qualified validation")
	}

	if c.transliteration {
		companyName = t.Transliterate(companyName)
		city = t.Transliterate(city)
		street = t.Transliterate(street)
	}

	req := &ValidationRequest{
		RequestingVATID: requestingVATID,
		RequestedVATID:  requestedVATID,
		CompanyName:     companyName,
		City:            city,
		Street:          street,
		PostalCode:      postalCode,
	}

	var resp ValidationResponse
	if err := c.doRequest(ctx, "POST", "/v1/abfrage", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ValidateVATWithRequest validates a VAT ID with a custom request.
func (c *Client) ValidateVATWithRequest(ctx context.Context, req *ValidationRequest) (*ValidationResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if err := validateRequestingVAT(req.RequestingVATID); err != nil {
		return nil, err
	}
	if err := validateRequestedVAT(req.RequestedVATID); err != nil {
		return nil, err
	}

	var resp ValidationResponse
	if err := c.doRequest(ctx, "POST", "/v1/abfrage", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func validateRequestingVAT(vat string) error {
	if vat == "" {
		return fmt.Errorf("requesting VAT ID is required")
	}
	if !strings.HasPrefix(vat, "DE") {
		return fmt.Errorf("requesting VAT ID must be German")
	}
	return nil
}

func validateRequestedVAT(vat string) error {
	if vat == "" {
		return fmt.Errorf("requested VAT ID is required")
	}
	if strings.HasPrefix(vat, "DE") {
		return fmt.Errorf("requested VAT ID cannot be German (DE)")
	}
	return nil
}
