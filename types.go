package evatr

import "time"

// ValidationRequest represents a VAT ID validation request.
type ValidationRequest struct {
	// Requesting German VAT ID (required)
	RequestingVATID string `json:"anfragendeUstid"`

	// VAT ID to validate (required)
	RequestedVATID string `json:"angefragteUstid"`

	// Company name (optional, required for qualified validation)
	CompanyName string `json:"firmenname,omitempty"`

	// Street address (optional)
	Street string `json:"strasse,omitempty"`

	// Postal code (optional)
	PostalCode string `json:"plz,omitempty"`

	// City (optional, required for qualified validation)
	City string `json:"ort,omitempty"`
}

// ValidationResponse represents a VAT ID validation response.
type ValidationResponse struct {
	// Technical ID for the validation request
	ID string `json:"id,omitempty"`

	// Timestamp of the request
	RequestTimestamp string `json:"anfrageZeitpunkt"`

	// Date from when the VAT ID is/was valid
	ValidFrom string `json:"gueltigAb,omitempty"`

	// Date until when the VAT ID was valid
	ValidUntil string `json:"gueltigBis,omitempty"`

	// Status code (e.g., "evatr-0000" for valid)
	Status string `json:"status"`

	// Company name verification result (A/B/C/D)
	CompanyNameResult VerificationResult `json:"ergFirmenname,omitempty"`

	// Street verification result (A/B/C/D)
	StreetResult VerificationResult `json:"ergStrasse,omitempty"`

	// Postal code verification result (A/B/C/D)
	PostalCodeResult VerificationResult `json:"ergPlz,omitempty"`

	// City verification result (A/B/C/D)
	CityResult VerificationResult `json:"ergOrt,omitempty"`
}

// GetRequestTimestamp parses the request timestamp.
func (v *ValidationResponse) GetRequestTimestamp() (time.Time, error) {
	return time.Parse(time.RFC3339, v.RequestTimestamp)
}

// GetValidFrom parses the valid-from date if present.
func (v *ValidationResponse) GetValidFrom() (time.Time, error) {
	if v.ValidFrom == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, v.ValidFrom)
}

// GetValidUntil parses the valid-until date if present.
func (v *ValidationResponse) GetValidUntil() (time.Time, error) {
	if v.ValidUntil == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, v.ValidUntil)
}

// IsValid returns whether the VAT ID is currently valid.
func (v *ValidationResponse) IsValid() bool {
	return v.Status == "evatr-0000" || v.Status == "evatr-2008"
}

// VerificationResult represents the result of a qualified validation field comparison.
type VerificationResult string

const (
	// Data matches the registered data
	VerificationMatch VerificationResult = "A"

	// Data does not match
	VerificationMismatch VerificationResult = "B"

	// Data was not requested
	VerificationNotRequested VerificationResult = "C"

	// Member state did not provide the data
	VerificationNotProvided VerificationResult = "D"
)

// StatusMessage represents a status message for an error code.
type StatusMessage struct {
	// Status code (e.g., "evatr-0000")
	Status string `json:"status"`

	// Category of the status
	Category string `json:"kategorie"`

	// Associated HTTP status code
	HTTPCode int `json:"httpcode"`

	// Field related to the status (if applicable)
	Field string `json:"feld,omitempty"`

	// Human-readable message
	Message string `json:"meldung"`
}

// EUMemberState represents an EU member state and its VIES availability.
type EUMemberState struct {
	// Two-letter country code
	Alpha2 string `json:"alpha2"`

	// Country name
	Name string `json:"name"`

	// Whether VIES system is available for this country
	Available bool `json:"verfuegbar"`
}
