package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hostwithquantum/go-evatr"
)

func main() {
	ctx := context.Background()

	// Create a new client
	client := evatr.NewClient()

	// Example 1: Simple VAT ID validation
	fmt.Println("=== Simple VAT ID Validation ===")
	simpleValidation(ctx, client)

	fmt.Println()

	// Example 2: Qualified validation with company data
	fmt.Println("=== Qualified VAT ID Validation ===")
	qualifiedValidation(ctx, client)

	fmt.Println()

	// Example 3: Get EU member states
	fmt.Println("=== EU Member States ===")
	getMemberStates(ctx, client)
}

func simpleValidation(ctx context.Context, client *evatr.Client) {
	// Use test values from the API documentation
	result, err := client.ValidateVAT(ctx, "DE123456789", "ATU12345678")
	if err != nil {
		if evatrErr, ok := err.(*evatr.Error); ok {
			fmt.Printf("Validation error: %s (HTTP %d)\n", evatrErr.Message, evatrErr.StatusCode)
			fmt.Printf("Status code: %s\n", evatrErr.Status)
			return
		}
		log.Fatal(err)
	}

	if result.IsValid() {
		fmt.Println("✓ VAT ID is valid")
		fmt.Printf("  Validated at: %s\n", result.RequestTimestamp)
		fmt.Printf("  Status: %s\n", result.Status)
	} else {
		fmt.Printf("✗ VAT ID is not valid: %s\n", result.Status)
		if result.ValidFrom != "" {
			fmt.Printf("  Valid from: %s\n", result.ValidFrom)
		}
		if result.ValidUntil != "" {
			fmt.Printf("  Valid until: %s\n", result.ValidUntil)
		}
	}
}

func qualifiedValidation(ctx context.Context, client *evatr.Client) {
	// Use test values from the API documentation
	result, err := client.ValidateVATQualified(
		ctx,
		"DE123456789",             // Your German VAT ID
		"ATU12345678",             // VAT ID to validate
		"Musterhaus GmbH & Co KG", // Company name
		"Musterort",               // City
		"Musterstrasse 22",        // Street
		"12345",                   // Postal code
	)

	if err != nil {
		if evatrErr, ok := err.(*evatr.Error); ok {
			fmt.Printf("Validation error: %s (HTTP %d)\n", evatrErr.Message, evatrErr.StatusCode)
			return
		}
		log.Fatal(err)
	}

	if result.IsValid() {
		fmt.Println("✓ VAT ID is valid")
		fmt.Printf("  Status: %s\n", result.Status)
		fmt.Println("\nCompany data verification:")
		printVerification("  Company name", result.CompanyNameResult)
		printVerification("  Street", result.StreetResult)
		printVerification("  Postal code", result.PostalCodeResult)
		printVerification("  City", result.CityResult)
	} else {
		fmt.Printf("✗ VAT ID is not valid: %s\n", result.Status)
	}
}

func getMemberStates(ctx context.Context, client *evatr.Client) {
	states, err := client.GetEUMemberStates(ctx)
	if err != nil {
		log.Fatal(err)
	}

	available := 0
	for _, state := range states {
		if state.Available {
			available++
		}
	}

	fmt.Printf("Found %d EU member states (%d available)\n", len(states), available)
	fmt.Println("\nSample member states:")

	// Show first 5 states as an example
	count := 0
	for _, state := range states {
		if count >= 5 {
			break
		}
		status := "✓ available"
		if !state.Available {
			status = "✗ unavailable"
		}
		fmt.Printf("  %s (%s): %s\n", state.Name, state.Alpha2, status)
		count++
	}

	if len(states) > 5 {
		fmt.Printf("  ... and %d more\n", len(states)-5)
	}
}

func printVerification(field, result string) {
	if result == "" {
		return
	}

	var msg string
	switch evatr.VerificationResult(result) {
	case evatr.VerificationMatch:
		msg = "✓ matches"
	case evatr.VerificationMismatch:
		msg = "✗ does not match"
	case evatr.VerificationNotRequested:
		msg = "- not requested"
	case evatr.VerificationNotProvided:
		msg = "- not provided by member state"
	default:
		msg = "unknown"
	}

	fmt.Printf("%s: %s\n", field, msg)
}
