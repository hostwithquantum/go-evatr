package evatr

import (
	"context"
	"fmt"
	"strings"
)

// ValidateVAT validates a VAT ID without company data verification.
func (c *Client) ValidateVAT(ctx context.Context, requestingVATID, requestedVATID string) (*ValidationResponse, error) {
	if requestingVATID == "" {
		return nil, fmt.Errorf("requesting VAT ID is required")
	}
	if !strings.HasPrefix(requestingVATID, "DE") {
		return nil, fmt.Errorf("requesting VAT ID must be German")
	}
	if requestedVATID == "" {
		return nil, fmt.Errorf("requested VAT ID is required")
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
	if requestingVATID == "" {
		return nil, fmt.Errorf("requesting VAT ID is required")
	}
	if requestedVATID == "" {
		return nil, fmt.Errorf("requested VAT ID is required")
	}
	if companyName == "" {
		return nil, fmt.Errorf("company name is required for qualified validation")
	}
	if city == "" {
		return nil, fmt.Errorf("city is required for qualified validation")
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
	if req.RequestingVATID == "" {
		return nil, fmt.Errorf("requesting VAT ID is required")
	}
	if req.RequestedVATID == "" {
		return nil, fmt.Errorf("requested VAT ID is required")
	}

	var resp ValidationResponse
	if err := c.doRequest(ctx, "POST", "/v1/abfrage", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
